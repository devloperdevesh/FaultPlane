"use client";

import { useEffect, useState } from "react";
import { getMetrics } from "@/lib/api";
import type { DashboardMetrics } from "@/lib/types";

export default function RuntimeHealth() {
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
    ["CPU", "—"],
    ["Memory", metrics ? `${metrics.memory.toFixed(1)} MB` : "—"],
    ["Network", "—"],
    ["Goroutines", "—"],
    ["Sockets", "—"],
  ];

  return (
    <div className="grid gap-4 md:grid-cols-5">
      {data.map(([name, value]) => (
        <div
          key={name}
          className="rounded-xl border border-white/10 bg-zinc-950 p-5"
        >
          <p className="text-xs text-zinc-500">{name}</p>

          <p className="mt-3 font-mono text-xl text-white">
            {value}
          </p>
        </div>
      ))}

      {error && (
        <div className="md:col-span-5 rounded-xl border border-white/10 bg-zinc-900/60 p-3 text-xs text-zinc-500">
          Runtime metrics unavailable.
        </div>
      )}
    </div>
  );
}
