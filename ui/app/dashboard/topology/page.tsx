"use client";

import { useEffect, useState } from "react";

import TopologyGraph from "@/components/topology/TopologyGraph";
import SocketMigration from "@/components/topology/SocketMigration";
import { getTopology } from "@/lib/api";
import type { TopologySnapshot } from "@/lib/api";

export default function TopologyPage() {
  const [topology, setTopology] = useState<TopologySnapshot | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    async function loadTopology() {
      try {
        const data = await getTopology();

        if (mounted) {
          setTopology(data);
          setError(null);
        }
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error
              ? err.message
              : "Failed to load topology",
          );
        }
      }
    }

    void loadTopology();

    const interval = window.setInterval(loadTopology, 5000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  return (
    <main className="min-h-screen bg-black p-6 text-white">
      <div className="mx-auto max-w-7xl space-y-6">
        <div>
          <h1 className="text-2xl font-semibold">
            Network Topology
          </h1>

          <p className="mt-1 text-sm text-zinc-500">
            Live runtime infrastructure and connection state
          </p>
        </div>

        {error && (
          <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-400">
            {error}
          </div>
        )}

        <div className="grid gap-6 lg:grid-cols-2">
          <TopologyGraph />

          <div className="space-y-4 rounded-xl border border-zinc-800 bg-zinc-950 p-6">
            <div>
              <h2 className="text-sm font-semibold text-white">
                Live Topology State
              </h2>

              <p className="mt-1 text-xs text-zinc-500">
                Data received from /api/topology
              </p>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="rounded-lg bg-zinc-900 p-4">
                <p className="text-xs text-zinc-500">Nodes</p>
                <p className="mt-1 text-xl font-semibold">
                  {topology?.nodes.length ?? 0}
                </p>
              </div>

              <div className="rounded-lg bg-zinc-900 p-4">
                <p className="text-xs text-zinc-500">Connections</p>
                <p className="mt-1 text-xl font-semibold">
                  {topology?.connections.length ?? 0}
                </p>
              </div>
            </div>

            <div className="rounded-lg bg-zinc-900 p-4">
              <p className="text-xs text-zinc-500">
                Last Updated
              </p>

              <p className="mt-1 text-sm text-zinc-300">
                {topology?.updated_at
                  ? new Date(topology.updated_at).toLocaleString()
                  : "Waiting for API..."}
              </p>
            </div>

            <div>
              <p className="mb-3 text-xs uppercase tracking-wider text-zinc-500">
                Connections
              </p>

              <div className="space-y-2">
                {topology?.connections.map((connection) => (
                  <div
                    key={connection.id}
                    className="rounded-lg border border-zinc-800 bg-black p-3"
                  >
                    <div className="flex items-center justify-between">
                      <span className="text-sm text-zinc-200">
                        {connection.source} → {connection.target}
                      </span>

                      <span className="text-xs text-emerald-400">
                        {connection.status}
                      </span>
                    </div>

                    <p className="mt-1 text-xs text-zinc-500">
                      {connection.type}
                    </p>
                  </div>
                ))}

                {!topology && (
                  <p className="text-sm text-zinc-600">
                    Loading topology...
                  </p>
                )}
              </div>
            </div>
          </div>
        </div>

        <SocketMigration />
      </div>
    </main>
  );
}
