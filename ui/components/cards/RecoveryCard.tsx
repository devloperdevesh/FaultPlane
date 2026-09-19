"use client";

import { useEffect, useState } from "react";
import { getMetrics } from "@/lib/api";
import type { DashboardMetrics } from "@/lib/types";

export default function RecoveryCard() {
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
      <p className="text-xs uppercase text-zinc-500">Recovery Events</p>

      <h3 className="mt-2 text-2xl font-semibold text-white">
        {metrics ? metrics.recoveries.toLocaleString() : "Unavailable"}
      </h3>

      <p className="mt-2 text-xs text-zinc-500">
        Recorded runtime recovery events
      </p>
    </div>
  );
}
