import type { TopologyData } from "./load";

// Explicit, deterministic replay fixtures. Never loaded into live API state.
export function demoSnapshot(step: number): TopologyData {
    const failed = step === 1;
    return {
        topology: {
            updated_at: `2026-01-01T00:00:0${step * 2}Z`,
            nodes: [
                {
                    id: "gateway",
                    name: "Demo gateway",
                    type: "gateway",
                    status: "healthy",
                },
                {
                    id: "worker-a",
                    name: "Demo worker A",
                    type: "worker",
                    status: failed ? "failed" : "active",
                    region: "demo-region-a",
                },
                {
                    id: "worker-b",
                    name: "Demo worker B",
                    type: "worker",
                    status: "active",
                    region: "demo-region-b",
                },
                {
                    id: "checkpoint",
                    name: "Demo checkpoint",
                    type: "checkpoint",
                    status: "healthy",
                    region: "demo-region-a",
                },
            ],
            connections: [
                {
                    id: "ga",
                    source: "gateway",
                    target: "worker-a",
                    type: "runtime",
                    status: "active",
                },
                {
                    id: "gb",
                    source: "gateway",
                    target: "worker-b",
                    type: "runtime",
                    status: "active",
                },
                {
                    id: "ac",
                    source: "worker-a",
                    target: "checkpoint",
                    type: "checkpoint-sync",
                    status: "active",
                },
            ],
        },
        workers: [
            {
                id: "worker-a",
                status: failed ? "failed" : "active",
                cpu: 12,
                memory: 16777216,
            },
            { id: "worker-b", status: "active", cpu: 9, memory: 8388608 },
        ],
        events:
            step === 0
                ? []
                : [
                      {
                          type: "worker_failed",
                          timestamp: "2026-01-01T00:00:02Z",
                          metadata: {
                              node_id: "worker-a",
                              region: "demo-region-a",
                              source: "demo",
                          },
                      },
                      ...(step === 2
                          ? [
                                {
                                    type: "worker_recovered",
                                    timestamp: "2026-01-01T00:00:04Z",
                                    metadata: {
                                        node_id: "worker-a",
                                        region: "demo-region-a",
                                        source: "demo",
                                    },
                                },
                            ]
                          : []),
                  ],
    };
}
