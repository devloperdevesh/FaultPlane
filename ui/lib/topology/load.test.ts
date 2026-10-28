import { afterEach, describe, expect, it, vi } from "vitest";
import { demoSnapshot } from "./demo";
import { loadTopologyData } from "./load";

afterEach(() => vi.unstubAllGlobals());

describe("topology API integration", () => {
    it("loads all three real endpoints with the same abort signal and no cache", async () => {
        const data = demoSnapshot(1);
        const fetcher = vi.fn(
            async (url: string) =>
                new Response(
                    JSON.stringify(
                        url.endsWith("/api/topology")
                            ? data.topology
                            : url.endsWith("/api/workers")
                              ? data.workers
                              : { events: data.events },
                    ),
                    { status: 200 },
                ),
        );
        vi.stubGlobal("fetch", fetcher);
        const controller = new AbortController();
        expect(await loadTopologyData(controller.signal)).toEqual(data);
        expect(
            fetcher.mock.calls.map(([url]) => new URL(url).pathname).sort(),
        ).toEqual(["/api/telemetry", "/api/topology", "/api/workers"]);
        for (const call of vi.mocked(fetch).mock.calls) {
            expect(call[1]?.signal).toBe(controller.signal);
            expect(call[1]?.cache).toBe("no-store");
        }
    });
    it("rejects HTTP failures rather than converting missing data into healthy metrics", async () => {
        vi.stubGlobal(
            "fetch",
            vi.fn(
                async () =>
                    new Response("", {
                        status: 503,
                        statusText: "Unavailable",
                    }),
            ),
        );
        await expect(
            loadTopologyData(new AbortController().signal),
        ).rejects.toThrow("503 Unavailable");
    });
    it("rejects malformed successful responses", async () => {
        vi.stubGlobal(
            "fetch",
            vi.fn(
                async () =>
                    new Response(JSON.stringify({ bad: "schema" }), {
                        status: 200,
                    }),
            ),
        );
        await expect(
            loadTopologyData(new AbortController().signal),
        ).rejects.toThrow("Invalid");
    });
});
