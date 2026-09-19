"use client";

import { useEffect, useState } from "react";
import { getTopology, type TopologyNode } from "@/lib/api";

const colors: Record<string, string> = {
  healthy: "border-emerald-500 bg-emerald-500/10 text-emerald-400",
  degraded: "border-orange-500 bg-orange-500/10 text-orange-400",
  failed: "border-red-500 bg-red-500/10 text-red-400",
};

export default function BlastGraph() {
  const [nodes, setNodes] = useState<TopologyNode[]>([]);
  const [error, setError] = useState(false);

  useEffect(() => {
    let active = true;

    const load = async () => {
      try {
        const response = await getTopology();
        if (active) {
          setNodes(response.nodes);
          setError(false);
        }
      } catch {
        if (active) setError(true);
      }
    };

    load();
    const interval = setInterval(load, 2000);

    return () => {
      active = false;
      clearInterval(interval);
    };
  }, []);

  return (
    <div className="rounded-xl border border-zinc-800 bg-black/40 p-6">
      <h3 className="mb-6 text-sm font-semibold text-white">
        Failure Dependency Graph
      </h3>

      {nodes.length === 0 ? (
        <div className="rounded-xl bg-zinc-900 p-4 text-sm text-zinc-500">
          {error
            ? "Topology API unavailable."
            : "No topology nodes available."}
        </div>
      ) : (
        <div className="flex flex-col items-center gap-6">
          {nodes.map((node) => {
            const status = node.status?.toLowerCase() || "unknown";

            return (
              <div
                key={node.id}
                className={`w-40 rounded-xl border px-4 py-3 text-center transition hover:scale-105 ${
                  colors[status] ??
                  "border-zinc-700 bg-zinc-900 text-zinc-400"
                }`}
              >
                <p className="text-sm font-medium">{node.name}</p>
                <p className="mt-1 text-xs capitalize">{status}</p>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
