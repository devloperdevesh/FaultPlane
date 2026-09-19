"use client";

import { useEffect, useState } from "react";
import { getMetrics } from "@/lib/api";
import type { DashboardMetrics } from "@/lib/types";

export default function LatencyChart() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);

  useEffect(() => {
    let active = true;

    const load = async () => {
      try {
        const data = await getMetrics();
        if (active) setMetrics(data);
      } catch {
        if (active) setMetrics(null);
      }
    };

    void load();
    const interval = window.setInterval(load, 2000);

    return () => {
      active = false;
      window.clearInterval(interval);
    };
  }, []);

  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-950 p-5">
      <h3 className="mb-4 text-sm font-medium text-white">
        Runtime Latency
      </h3>
      <div className="flex h-[260px] items-center justify-center">
        <div className="text-center">
          <p className="text-3xl font-semibold text-white">
            {metrics ? `${metrics.latency.toFixed(2)} ms` : "Unavailable"}
          </p>
          <p className="mt-2 text-xs text-zinc-500">
            Live average latency from /api/metrics
          </p>
        </div>
      </div>
    </div>
  );
}
