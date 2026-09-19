"use client";

import { useEffect, useState } from "react";
import { getMetrics } from "@/lib/api";
import MemoryStats from "./MemoryStats";

export default function MemoryGrid() {
  const [memory, setMemory] = useState<number | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    let active = true;

    const load = async () => {
      try {
        const metrics = await getMetrics();

        if (active) {
          setMemory(metrics.memory);
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
    <section className="rounded-2xl border border-zinc-800 bg-zinc-950 p-6 shadow-xl shadow-black/20">
      <div className="mb-8 flex items-start justify-between">
        <div>
          <h2 className="text-base font-semibold text-white">
            Runtime Memory
          </h2>
          <p className="mt-1 text-xs text-zinc-500">
            Live runtime memory usage reported by FaultPlane
          </p>
        </div>

        <div className="rounded-full border border-zinc-700 bg-zinc-900 px-3 py-1">
          <span className="text-xs font-mono text-zinc-400">
            {error ? "UNAVAILABLE" : "LIVE"}
          </span>
        </div>
      </div>

      <div className="rounded-xl border border-zinc-800 bg-black/40 p-5">
        <p className="mb-3 text-xs uppercase tracking-wider text-zinc-500">
          Runtime Memory
        </p>

        <div className="text-2xl font-mono text-white">
          {memory === null ? "—" : `${memory.toFixed(2)} MB`}
        </div>

        <p className="mt-2 text-xs text-zinc-600">
          Source: /api/metrics
        </p>
      </div>

      <div className="mt-6 grid gap-4 md:grid-cols-3">
        <MemoryStats
          title="Runtime Memory"
          value={memory === null ? "—" : `${memory.toFixed(2)} MB`}
        />
        <MemoryStats title="Snapshot Size" value="—" />
        <MemoryStats title="Migration Latency" value="—" />
      </div>
    </section>
  );
}
