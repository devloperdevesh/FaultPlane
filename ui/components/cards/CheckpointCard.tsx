"use client";

import { useEffect, useState } from "react";
import { getMetrics } from "../../lib/api";
import type { DashboardMetrics } from "../../lib/types";

export default function CheckpointCard() {
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
    const timer = window.setInterval(() => void load(), 2000);

    return () => {
      active = false;
      window.clearInterval(timer);
    };
  }, []);

  return (
    <div className="rounded-xl border border-white/10 bg-zinc-950 p-5">
      <p className="text-xs uppercase tracking-widest text-zinc-500">
        Checkpoints
      </p>
      <p className="mt-2 text-3xl font-semibold text-white">
        {metrics ? metrics.checkpoints.toLocaleString() : "Unavailable"}
      </p>
      <p className="mt-2 text-xs text-zinc-500">
        Runtime checkpoints recorded
      </p>
    </div>
  );
}
