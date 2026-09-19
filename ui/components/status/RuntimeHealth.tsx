"use client";

import { useEffect, useState } from "react";
import { getMetrics } from "@/lib/api";
import type { DashboardMetrics } from "@/lib/types";

export default function RuntimeHealth() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);

  useEffect(() => {
    let active = true;

    const loadMetrics = async () => {
      try {
        const data = await getMetrics();
        if (active) setMetrics(data);
      } catch {
        if (active) setMetrics(null);
      }
    };

    void loadMetrics();
    const interval = window.setInterval(loadMetrics, 2000);

    return () => {
      active = false;
      window.clearInterval(interval);
    };
  }, []);

  if (!metrics) {
    return (
      <div className="rounded-xl border border-white/10 bg-zinc-950 p-5 text-sm text-zinc-500">
        Runtime metrics unavailable.
      </div>
    );
  }

  const data = [
    ["CPU", `${metrics.cpu.toFixed(2)}%`],
    ["Memory", `${metrics.memory.toFixed(1)} MB`],
    ["Workers", metrics.workers.toLocaleString()],
    ["Requests", metrics.requests.toLocaleString()],
  ];

  return (
    <div className="rounded-xl border border-white/10 bg-zinc-950 p-5">
      <h3 className="mb-4 text-sm font-semibold text-white">Runtime Health</h3>

      <div className="grid grid-cols-2 gap-4">
        {data.map(([label, value]) => (
          <div key={label} className="rounded-lg bg-zinc-900/70 p-4">
            <p className="text-xs text-zinc-500">{label}</p>
            <p className="mt-1 text-lg font-medium text-white">{value}</p>
          </div>
        ))}
      </div>

      <p className="mt-4 text-xs text-zinc-600">
        Source: /api/metrics
      </p>
    </div>
  );
}
