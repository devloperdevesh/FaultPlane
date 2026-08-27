"use client";

import { useEffect, useState } from "react";
import WorkerStatus from "./WorkerStatus";
import WorkerHealth from "./WorkerHealth";
import { getWorkers, type RuntimeWorker } from "@/lib/api";

export default function WorkerTable() {
  const [workers, setWorkers] = useState<RuntimeWorker[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;

    const load = async () => {
      try {
        const data = await getWorkers();

        if (mounted) {
          setWorkers(data);
          setError(null);
          setLoading(false);
        }
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load runtime workers",
          );
          setLoading(false);
        }
      }
    };

    void load();

    const interval = window.setInterval(() => {
      void load();
    }, 2000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  return (
    <div className="overflow-hidden rounded-2xl border border-white/10 bg-zinc-950">
      <div className="border-b border-white/10 px-6 py-5">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-sm font-semibold text-white">
              Runtime workers
            </h2>
            <p className="mt-1 text-xs text-zinc-500">
              Live worker inventory reported by FaultPlane
            </p>
          </div>

          <span className="rounded-full border border-white/10 px-3 py-1 text-[11px] text-zinc-400">
            {workers.length} registered
          </span>
        </div>
      </div>

      {loading ? (
        <div className="px-6 py-12 text-center text-sm text-zinc-500">
          Loading runtime workers...
        </div>
      ) : error ? (
        <div className="px-6 py-12 text-center">
          <p className="text-sm text-red-400">
            Runtime workers unavailable
          </p>
          <p className="mt-2 text-xs text-zinc-600">{error}</p>
        </div>
      ) : workers.length === 0 ? (
        <div className="px-6 py-12 text-center">
          <p className="text-sm text-zinc-400">
            No workers registered
          </p>
          <p className="mt-2 text-xs text-zinc-600">
            The runtime has not reported any active workers yet.
          </p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="border-b border-white/10 text-left text-[11px] uppercase tracking-wider text-zinc-600">
              <tr>
                <th className="px-6 py-4">Worker</th>
                <th className="px-4 py-4">Status</th>
                <th className="px-4 py-4">CPU</th>
                <th className="px-4 py-4">Memory</th>
              </tr>
            </thead>

            <tbody>
              {workers.map((worker) => (
                <tr
                  key={worker.id}
                  className="border-b border-white/[0.06] transition hover:bg-white/[0.03]"
                >
                  <td className="px-6 py-4 font-mono text-white">
                    {worker.id}
                  </td>

                  <td className="px-4 py-4">
                    <WorkerStatus status={worker.status as "ACTIVE" | "RECOVERING" | "FAILED"} />
                  </td>

                  <td className="px-4 py-4 text-zinc-300">
                    {worker.cpu > 0
                      ? `${worker.cpu.toFixed(1)}%`
                      : "—"}
                  </td>

                  <td className="px-4 py-4 text-zinc-300">
                    {worker.memory > 0
                      ? `${worker.memory} B`
                      : "—"}
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
