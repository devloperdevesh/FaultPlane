// @vitest-environment jsdom
import { act, createElement } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import TopologyDashboard from "../../components/topology/TopologyDashboard";
import TopologyPage from "../../app/dashboard/topology/page";
import { demoSnapshot } from "./demo";
import { loadTopologyData } from "./load";

vi.mock("./load", async (original) => ({
    ...(await original<typeof import("./load")>()),
    loadTopologyData: vi.fn(),
}));
const load = vi.mocked(loadTopologyData);
Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
let container: HTMLDivElement;
let root: Root;
const text = () => container.textContent ?? "";
const button = (name: string) => {
    const result = [...container.querySelectorAll("button")].find((item) =>
        item.textContent?.includes(name),
    );
    if (!result) throw new Error(`Button not found: ${name}`);
    return result;
};
async function click(name: string) {
    await act(async () => {
        button(name).click();
    });
}

beforeEach(() => {
    vi.useFakeTimers();
    load.mockReset();
    container = document.createElement("div");
    document.body.append(container);
    root = createRoot(container);
});
afterEach(async () => {
    await act(async () => root.unmount());
    container.remove();
    vi.useRealTimers();
});

describe("topology HUD", () => {
    it("replays a scoped failure and recovery without ever querying live APIs", async () => {
        await act(async () =>
            root.render(createElement(TopologyDashboard, { demo: true })),
        );
        expect(text()).toContain("DEMO — fixture replay");
        expect(button("Demo worker A").textContent).toContain(
            "Observed failure",
        );
        expect(button("Demo checkpoint").textContent).toContain(
            "Potential exposure",
        );
        expect(button("Demo worker B").textContent).toContain("healthy");
        await click("Demo worker A");
        expect(
            container.querySelector("#node-inspector")?.textContent,
        ).toContain("12.0% / 16.0 MiB");
        expect(button("Demo worker A").getAttribute("aria-pressed")).toBe(
            "true",
        );
        await click("3. Recovered");
        expect(button("Demo worker A").textContent).toContain("healthy");
        expect(button("Demo checkpoint").textContent).not.toContain(
            "Potential exposure",
        );
        expect(text()).toContain("Recovery observed");
        expect(text()).toContain("2026-01-01T00:00:04.000Z");
        await click("1. Healthy");
        expect(text()).toContain("No failure or recovery events recorded");
        expect(load).not.toHaveBeenCalled();
    });
    it("displays loading then a disconnected state and retries successfully", async () => {
        let reject!: (reason: Error) => void;
        load.mockImplementationOnce(
            () =>
                new Promise((_, failure) => {
                    reject = failure;
                }),
        );
        await act(async () => root.render(createElement(TopologyDashboard)));
        expect(text()).toContain("Loading runtime telemetry");
        await act(async () => reject(new Error("Network unavailable")));
        expect(text()).toContain("Disconnected / API error");
        expect(text()).toContain("Network unavailable");
        load.mockResolvedValue(demoSnapshot(0));
        await click("Retry now");
        expect(text()).toContain("Connected · polling");
        expect(text()).not.toContain("Disconnected / API error");
    });
    it("retains the last good graph on error and replaces it on reconnect", async () => {
        load.mockResolvedValueOnce(demoSnapshot(1))
            .mockRejectedValueOnce(new Error("503 Service unavailable"))
            .mockResolvedValue(demoSnapshot(2));
        await act(async () => root.render(createElement(TopologyDashboard)));
        await act(async () => {
            await vi.advanceTimersByTimeAsync(5000);
        });
        expect(text()).toContain("stale snapshot");
        expect(button("Demo worker A").textContent).toContain(
            "Observed failure",
        );
        await act(async () => {
            await vi.advanceTimersByTimeAsync(5000);
        });
        expect(text()).not.toContain("stale snapshot");
        expect(button("Demo worker A").textContent).toContain("healthy");
    });
    it("renders a successful empty snapshot without inventing workers", async () => {
        load.mockResolvedValue({
            topology: {
                nodes: [],
                connections: [],
                updated_at: "2026-01-01T00:00:00Z",
            },
            events: [],
            workers: [],
        });
        await act(async () => root.render(createElement(TopologyDashboard)));
        expect(text()).toContain("No topology nodes reported");
        expect(text()).toContain("No failure or recovery events recorded");
        expect(text()).not.toContain("Demo worker");
    });
    it("aborts hung fetches, never overlaps polls and cancels timers on unmount", async () => {
        load.mockImplementation(
            (signal) =>
                new Promise((_, reject) =>
                    signal.addEventListener(
                        "abort",
                        () => reject(new Error("Request timed out")),
                        { once: true },
                    ),
                ),
        );
        await act(async () => root.render(createElement(TopologyDashboard)));
        await act(async () => {
            await vi.advanceTimersByTimeAsync(9999);
        });
        expect(load).toHaveBeenCalledTimes(1);
        await act(async () => {
            await vi.advanceTimersByTimeAsync(1);
        });
        expect(text()).toContain("Request timed out");
        await act(async () => {
            await vi.advanceTimersByTimeAsync(5000);
        });
        expect(load).toHaveBeenCalledTimes(2);
        await act(async () => root.unmount());
        expect(load.mock.calls[1][0].aborted).toBe(true);
        await act(async () => {
            await vi.advanceTimersByTimeAsync(20000);
        });
        expect(load).toHaveBeenCalledTimes(2);
        root = createRoot(container);
    });
    it("isolates mode changes and aborts in-flight live data when opening demo", async () => {
        load.mockImplementation(() => new Promise(() => {}));
        await act(async () => root.render(createElement(TopologyPage)));
        const signal = load.mock.calls[0][0];
        await click("Open demo replay");
        expect(signal.aborted).toBe(true);
        expect(text()).toContain("DEMO — fixture replay");
        load.mockResolvedValue({ ...demoSnapshot(0), events: [] });
        await click("Return to live telemetry");
        expect(text()).not.toContain("DEMO — fixture replay");
        expect(text()).toContain("Connected · polling");
    });
});
