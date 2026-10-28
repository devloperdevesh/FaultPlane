import type { TelemetryEvent, TopologyNode, TopologySnapshot } from "../types";

export type Health = "healthy" | "failed" | "unknown";
export type EventKind = "failure" | "recovery";

export function health(status: string): Health {
    switch (status.trim().toLowerCase()) {
        case "failed":
        case "unhealthy":
        case "offline":
        case "disconnected":
        case "dead":
        case "error":
            return "failed";
        case "active":
        case "healthy":
        case "running":
        case "ready":
        case "recovered":
            return "healthy";
        default:
            return "unknown";
    }
}

export function eventKind(event: TelemetryEvent): EventKind | null {
    switch (event.type) {
        case "worker_failed":
        case "connection_failed":
        case "tcp_connection_failed":
        case "recovery_failed":
            return "failure";
        case "worker_recovered":
        case "connection_recovered":
        case "recovery_success":
            return "recovery";
        default:
            return null;
    }
}

export function regionOf(node: TopologyNode): string {
    return node.region?.trim() || "Region not reported";
}

export function eventTargets(
    event: TelemetryEvent,
    snapshot: TopologySnapshot,
): string[] {
    // An explicit but unknown node must not fall back to affecting an entire region.
    const id = event.metadata?.node_id || event.metadata?.worker_id;
    if (id) return snapshot.nodes.some((node) => node.id === id) ? [id] : [];
    const region = event.metadata?.region?.trim();
    return region
        ? snapshot.nodes
              .filter((node) => node.region?.trim() === region)
              .map((node) => node.id)
        : [];
}

export function incidentEvents(events: TelemetryEvent[]): TelemetryEvent[] {
    return events
        .filter(
            (event) =>
                eventKind(event) &&
                Number.isFinite(Date.parse(event.timestamp)),
        )
        .sort((a, b) => Date.parse(a.timestamp) - Date.parse(b.timestamp));
}

export function buildImpact(
    snapshot: TopologySnapshot,
    events: TelemetryEvent[],
) {
    const latest = new Map<string, TelemetryEvent>();
    const timeline = incidentEvents(events);
    for (const event of timeline) {
        for (const id of eventTargets(event, snapshot)) latest.set(id, event);
    }
    const failed = new Set<string>();
    for (const node of snapshot.nodes) {
        if (
            health(node.status) === "failed" ||
            (latest.has(node.id) &&
                eventKind(latest.get(node.id)!) === "failure")
        ) {
            failed.add(node.id);
        }
    }
    // Follow only reported directed edges. Reachability is potential exposure,
    // never evidence that those dependencies actually failed.
    const known = new Set(snapshot.nodes.map((node) => node.id));
    const adjacency = new Map<string, string[]>();
    for (const edge of snapshot.connections) {
        if (known.has(edge.source) && known.has(edge.target)) {
            adjacency.set(edge.source, [
                ...(adjacency.get(edge.source) ?? []),
                edge.target,
            ]);
        }
    }
    const reachable = new Set(failed);
    const queue = [...failed];
    for (let i = 0; i < queue.length; i++) {
        for (const target of adjacency.get(queue[i]) ?? []) {
            if (!reachable.has(target)) {
                reachable.add(target);
                queue.push(target);
            }
        }
    }
    const atRisk = new Set([...reachable].filter((id) => !failed.has(id)));
    const regions = new Map<string, TopologyNode[]>();
    for (const node of snapshot.nodes) {
        const name = regionOf(node);
        regions.set(name, [...(regions.get(name) ?? []), node]);
    }
    return { failed, atRisk, latest, regions, timeline };
}
