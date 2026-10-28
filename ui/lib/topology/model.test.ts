import { describe, expect, it } from "vitest";
import { demoSnapshot } from "./demo";
import {
    buildImpact,
    eventTargets,
    health,
    incidentEvents,
    regionOf,
} from "./model";
import { parseTopologyData } from "./load";
import type { TelemetryEvent } from "../types";

const event = (
    type: string,
    node_id: string,
    timestamp = "2026-01-01T00:00:02Z",
): TelemetryEvent => ({ type, timestamp, metadata: { node_id } });

describe("regional blast radius", () => {
    it("follows a recorded failure through outgoing edges, without declaring downstream failures", () => {
        const data = demoSnapshot(1);
        const impact = buildImpact(data.topology, data.events);
        expect([...impact.failed]).toEqual(["worker-a"]);
        expect([...impact.atRisk]).toEqual(["checkpoint"]);
        expect(impact.failed.has("worker-b")).toBe(false);
        expect(impact.atRisk.has("gateway")).toBe(false);
        expect(impact.regions.size).toBe(3);
    });
    it("clears exposure after a scoped recovery, regardless of input order", () => {
        const data = demoSnapshot(2);
        const impact = buildImpact(data.topology, [...data.events].reverse());
        expect(impact.failed.size).toBe(0);
        expect(impact.atRisk.size).toBe(0);
        expect(impact.latest.get("worker-a")?.type).toBe("worker_recovered");
        expect(impact.timeline).toHaveLength(2);
    });
    it("retains a failed snapshot even after an older recovery event", () => {
        const data = demoSnapshot(1);
        expect(
            buildImpact(data.topology, [
                event("worker_recovered", "worker-a"),
            ]).failed.has("worker-a"),
        ).toBe(true);
    });
    it("does not clear another region when a worker recovers", () => {
        const data = demoSnapshot(0);
        const events = [
            event("worker_failed", "worker-a"),
            event("worker_failed", "worker-b"),
            event("worker_recovered", "worker-a", "2026-01-01T00:00:04Z"),
        ];
        expect([...buildImpact(data.topology, events).failed]).toEqual([
            "worker-b",
        ]);
    });
    it("never expands an unknown node ID into its claimed region", () => {
        const data = demoSnapshot(0);
        const unknown = {
            ...event("worker_failed", "missing"),
            metadata: { node_id: "missing", region: "demo-region-a" },
        };
        expect(eventTargets(unknown, data.topology)).toEqual([]);
        expect(buildImpact(data.topology, [unknown]).timeline).toHaveLength(1);
        expect(buildImpact(data.topology, [unknown]).failed.size).toBe(0);
    });
    it("can correlate an explicitly regional event without inventing coordinates", () => {
        const data = demoSnapshot(0);
        const regional = {
            ...event("connection_failed", ""),
            metadata: { region: "demo-region-a" },
        };
        expect(eventTargets(regional, data.topology)).toEqual([
            "worker-a",
            "checkpoint",
        ]);
        expect(regionOf(data.topology.nodes[0])).toBe("Region not reported");
    });
    it("terminates on cyclic graphs and ignores dangling edges", () => {
        const data = demoSnapshot(1);
        data.topology.connections.push(
            {
                id: "cycle",
                source: "checkpoint",
                target: "worker-a",
                type: "runtime",
                status: "active",
            },
            {
                id: "dangling",
                source: "worker-a",
                target: "missing",
                type: "runtime",
                status: "active",
            },
        );
        expect([...buildImpact(data.topology, data.events).atRisk]).toEqual([
            "checkpoint",
        ]);
    });
    it("does not interpret generic HTTP activity, raw TCP observations or invalid timestamps as incidents", () => {
        expect(
            incidentEvents([
                event("request_complete", "worker-a"),
                event("tcp_set_state", "worker-a"),
                event("worker_failed", "worker-a", "invalid"),
            ]),
        ).toEqual([]);
        expect(health("unrecognised")).toBe("unknown");
        expect(health(" ACTIVE ")).toBe("healthy");
    });
});

describe("API boundary validation", () => {
    it("accepts Go null slices as empty and preserves valid data", () => {
        const data = demoSnapshot(1);
        expect(
            parseTopologyData(
                data.topology,
                { events: data.events },
                data.workers,
            ),
        ).toEqual(data);
        expect(
            parseTopologyData(
                {
                    nodes: null,
                    connections: null,
                    updated_at: data.topology.updated_at,
                },
                { events: null },
                null,
            ).topology.nodes,
        ).toEqual([]);
    });
    it.each([
        null,
        {},
        { nodes: "invalid" },
        { ...demoSnapshot(0).topology, nodes: [{ id: 3 }] },
        { ...demoSnapshot(0).topology, updated_at: "invalid" },
    ])("rejects malformed topology %j", (topology) => {
        expect(() => parseTopologyData(topology, { events: [] }, [])).toThrow();
    });
    it("rejects duplicate node IDs, invalid event metadata and invalid worker resources", () => {
        const data = demoSnapshot(0);
        expect(() =>
            parseTopologyData(
                {
                    ...data.topology,
                    nodes: [data.topology.nodes[0], data.topology.nodes[0]],
                },
                { events: [] },
                [],
            ),
        ).toThrow();
        expect(() =>
            parseTopologyData(
                data.topology,
                {
                    events: [
                        {
                            type: "worker_failed",
                            timestamp: "now",
                            metadata: { node_id: 3 },
                        },
                    ],
                },
                [],
            ),
        ).toThrow();
        expect(() =>
            parseTopologyData(data.topology, { events: [] }, [
                { id: "a", cpu: "12" },
            ]),
        ).toThrow();
    });
});
