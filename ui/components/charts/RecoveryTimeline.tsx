"use client";

import { useEffect, useState } from "react";
import { getMetrics } from "@/lib/api";
import type { DashboardMetrics } from "@/lib/types";

export default function RecoveryTimeline() {
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
      <h3 className="mb-4 text-sm text-white">Recovery Timeline</h3>

      <div className="flex min-h-[240px] items-center justify-center rounded-lg border border-white/5 bg-zinc-900/40 p-6 text-center">
        {metrics && metrics.recoveries > 0 ? (
          <div>
            <p className="text-2xl font-semibold text-white">
              {metrics.recoveries.toLocaleString()}
            </p>
            <p className="mt-2 text-xs text-zinc-500">
              Recorded recovery events. Detailed recovery timeline data is not
              currently exposed by the backend.
            </p>
          </div>
        ) : (
          <p className="text-sm text-zinc-500">
            No recovery timeline events available.
          </p>
        )}
      </div>
    </div>
  );
}
