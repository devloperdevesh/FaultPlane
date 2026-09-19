"use client";

import { useEffect, useState } from "react";
import { getHealth, getWorkers } from "@/lib/api";
import type { HealthResponse, RuntimeWorker } from "@/lib/types";

export default function RuntimeStatus() {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [workers, setWorkers] = useState<RuntimeWorker[]>([]);
  const [error, setError] = useState(false);

  useEffect(() => {
    let mounted = true;

    async function loadStatus() {
      try {
        const [healthData, workerData] = await Promise.all([
          getHealth(),
          getWorkers(),
        ]);

        if (mounted) {
          setHealth(healthData);
          setWorkers(workerData);
          setError(false);
        }
      } catch {
        if (mounted) {
          setError(true);
        }
      }
    }

    void loadStatus();

    const interval = window.setInterval(loadStatus, 2000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  const workerStatus =
    workers.length > 0
      ? `${workers.length} ${workers.length === 1 ? "Worker" : "Workers"}`
      : "No workers reported";

  const statusItems = [
    {
      name: "Gateway",
      status: health?.status ?? "—",
    },
    {
      name: "Storage",
      status: "Not reported",
    },
    {
      name: "Workers",
      status: workerStatus,
    },
  ];

  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-950 p-5">
      <h2 className="mb-5 text-sm font-semibold text-white">
        Runtime Status
      </h2>

      <div className="space-y-3">
        {statusItems.map((item) => (
          <div
            key={item.name}
            className="flex justify-between"
          >
            <span className="text-sm text-zinc-400">
              {item.name}
            </span>

            <span className="text-sm text-zinc-300">
              {item.status}
            </span>
          </div>
        ))}
      </div>

      {error && (
        <div className="mt-4 rounded-lg border border-white/10 bg-zinc-900/60 p-3 text-xs text-zinc-500">
          Runtime status unavailable.
        </div>
      )}
    </div>
  );
}
