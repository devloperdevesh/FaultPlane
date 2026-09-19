"use client";

import { useEffect, useState } from "react";
import { getWorkers } from "@/lib/api";
import type { RuntimeWorker } from "@/lib/types";

function formatMemory(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes < 0) return "—";
  if (bytes >= 1024 ** 3) return `${(bytes / 1024 ** 3).toFixed(1)} GB`;
  if (bytes >= 1024 ** 2) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  if (bytes >= 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${bytes} B`;
}

function statusClass(status: string): string {
  const normalized = status.toLowerCase();

  if (normalized === "healthy" || normalized === "running") {
    return "text-emerald-400";
  }

  if (normalized === "recovering" || normalized === "degraded") {
    return "text-amber-400";
  }

  if (normalized === "failed" || normalized === "offline") {
    return "text-red-400";
  }

  return "text-zinc-400";
}

export default function WorkersTable() {
  const [workers, setWorkers] = useState<RuntimeWorker[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    const load = async () => {
      try {
        const data = await getWorkers();

        if (!mounted) return;

        setWorkers(Array.isArray(data) ? data : []);
        setError(null);
      } catch (err) {
        if (!mounted) return;

        setError(err instanceof Error ? err.message : "Unable to load workers");
      } finally {
        if (mounted) setLoading(false);
      }
    };

    void load();

    const interval = window.setInterval(load, 2000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  return (
    <div className="overflow-hidden rounded-xl border border-white/10 bg-zinc-950/80">
      <div className="border-b border-white/10 px-5 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="font-medium text-white">Runtime Workers</h3>
            <p className="mt-1 text-xs text-zinc-500">
              Live worker state from the FaultPlane runtime API
            </p>
          </div>

          <span className="text-xs text-zinc-500">
            {loading ? "Loading…" : `${workers.length} workers`}
          </span>
        </div>
      </div>

      {error ? (
        <div className="px-5 py-8 text-sm text-red-400">
          <p>Worker API unavailable.</p>
          <p className="mt-1 text-xs text-zinc-500">{error}</p>
        </div>
      ) : workers.length === 0 && !loading ? (
        <div className="px-5 py-8 text-center text-sm text-zinc-500">
          No runtime workers reported by the API.
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-white/10 text-xs uppercase tracking-wider text-zinc-500">
              <tr>
                <th className="px-5 py-3 font-medium">Worker</th>
                <th className="px-5 py-3 font-medium">Status</th>
                <th className="px-5 py-3 font-medium">CPU</th>
                <th className="px-5 py-3 font-medium">Memory</th>
              </tr>
            </thead>

            <tbody>
              {workers.map((worker) => (
                <tr
                  key={worker.id}
                  className="border-b border-white/5 last:border-0"
                >
                  <td className="px-5 py-4 font-mono text-xs text-white">
                    {worker.id}
                  </td>

                  <td
                    className={`px-5 py-4 text-xs font-medium ${statusClass(
                      worker.status,
                    )}`}
                  >
                    {worker.status || "unknown"}
                  </td>

                  <td className="px-5 py-4 text-zinc-300">
                    {Number.isFinite(worker.cpu)
                      ? `${worker.cpu.toFixed(1)}%`
                      : "—"}
                  </td>

                  <td className="px-5 py-4 text-zinc-300">
                    {formatMemory(worker.memory)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}