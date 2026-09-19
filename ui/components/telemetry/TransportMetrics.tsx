"use client";

import { useEffect, useState } from "react";
import { getMetrics } from "@/lib/api";
import type { DashboardMetrics } from "@/lib/types";

export default function TransportMetrics() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    let mounted = true;

    async function loadMetrics() {
      try {
        const data = await getMetrics();

        if (mounted) {
          setMetrics(data);
          setError(false);
        }
      } catch {
        if (mounted) {
          setError(true);
        }
      }
    }

    void loadMetrics();

    const interval = window.setInterval(loadMetrics, 2000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  const data = [
    ["Requests", metrics ? metrics.requests.toLocaleString() : "—"],
    ["Latency", metrics ? `${metrics.latency.toFixed(2)} ms` : "—"],
    ["Connections", "—"],
    ["Dropped", "—"],
  ];

  return (
    <div className="grid gap-4 md:grid-cols-4">
      {data.map(([name, value]) => (
        <div
          key={name}
          className="rounded-xl border border-white/10 bg-zinc-950 p-4"
        >
          <p className="text-xs text-zinc-500">{name}</p>

          <p className="mt-2 font-mono text-white">
            {value}
          </p>
        </div>
      ))}

      {error && (
        <div className="md:col-span-4 rounded-xl border border-white/10 bg-zinc-900/60 p-3 text-xs text-zinc-500">
          Transport metrics unavailable.
        </div>
      )}
    </div>
  );
}
