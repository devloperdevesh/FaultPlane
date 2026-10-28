import { getTopology, getTelemetry, getWorkers } from "../api";
import type { TelemetryEvent, TopologySnapshot, RuntimeWorker } from "../types";

export interface TopologyData {
    topology: TopologySnapshot;
    events: TelemetryEvent[];
    workers: RuntimeWorker[];
}

function record(value: unknown): value is Record<string, unknown> {
    return typeof value === "object" && value !== null && !Array.isArray(value);
}
const text = (value: unknown): value is string => typeof value === "string";
const optionalText = (value: unknown) => value === undefined || text(value);

// Decode at the API boundary; Go nil slices are valid empty collections.
export function parseTopologyData(
    topology: unknown,
    telemetry: unknown,
    workers: unknown,
): TopologyData {
    if (!record(topology) || !record(telemetry))
        throw new Error("Invalid topology or telemetry response");
    const nodes = topology.nodes ?? [];
    const connections = topology.connections ?? [];
    const events = telemetry.events ?? [];
    const workerList = workers ?? [];
    if (
        !text(topology.updated_at) ||
        !Number.isFinite(Date.parse(topology.updated_at)) ||
        !Array.isArray(nodes) ||
        !Array.isArray(connections) ||
        !Array.isArray(events) ||
        !Array.isArray(workerList) ||
        !nodes.every(
            (node) =>
                record(node) &&
                text(node.id) &&
                node.id.length > 0 &&
                text(node.name) &&
                text(node.type) &&
                text(node.status) &&
                optionalText(node.region),
        ) ||
        !connections.every(
            (edge) =>
                record(edge) &&
                text(edge.id) &&
                text(edge.source) &&
                text(edge.target) &&
                text(edge.type) &&
                text(edge.status),
        ) ||
        !events.every(
            (event) =>
                record(event) &&
                text(event.type) &&
                text(event.timestamp) &&
                (event.metadata === undefined ||
                    (record(event.metadata) &&
                        Object.values(event.metadata).every(text))),
        ) ||
        !workerList.every(
            (worker) =>
                record(worker) &&
                text(worker.id) &&
                text(worker.status) &&
                typeof worker.cpu === "number" &&
                Number.isFinite(worker.cpu) &&
                typeof worker.memory === "number" &&
                Number.isFinite(worker.memory),
        ) ||
        new Set(nodes.map((node) => node.id)).size !== nodes.length
    ) {
        throw new Error(
            "Invalid topology, worker or telemetry data; the last valid snapshot is preserved",
        );
    }
    return {
        topology: { nodes, connections, updated_at: topology.updated_at },
        events,
        workers: workerList,
    };
}

export async function loadTopologyData(
    signal: AbortSignal,
): Promise<TopologyData> {
    const [topology, telemetry, workers] = await Promise.all([
        getTopology(signal),
        getTelemetry(signal),
        getWorkers(signal),
    ]);
    return parseTopologyData(topology, telemetry, workers);
}
