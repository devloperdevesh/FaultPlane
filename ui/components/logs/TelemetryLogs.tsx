"use client";

import { useEffect, useState } from "react";
import { getLogs, type TelemetryEvent } from "@/lib/api";

const colors: Record<string, string> = {
  INFO: "text-blue-400",
  SUCCESS: "text-emerald-400",
  WARN: "text-yellow-400",
  ERROR: "text-red-400",
};

export default function TelemetryLogs() {
  const [logs, setLogs] = useState<TelemetryEvent[]>([]);
  const [error, setError] = useState(false);

  useEffect(() => {
    let active = true;

    const load = async () => {
      try {
        const response = await getLogs();
        if (active) {
          setLogs(response.logs);
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
    <div className="rounded-2xl border border-white/10 bg-zinc-950 p-6">
      <div className="mb-5 flex justify-between">
        <h2 className="text-sm font-semibold text-white">
          Runtime Telemetry
        </h2>

        <span className="text-xs text-emerald-400">
          {error ? "UNAVAILABLE" : "LIVE"}
        </span>
      </div>

      <div className="space-y-3">
        {logs.length === 0 ? (
          <div className="rounded-xl bg-zinc-900 p-4 text-sm text-zinc-500">
            {error
              ? "Telemetry API unavailable."
              : "No telemetry events available."}
          </div>
        ) : (
          logs.map((log, index) => {
            const level =
              log.metadata?.level?.toUpperCase() ||
              log.type?.toUpperCase() ||
              "INFO";

            return (
              <div
                key={`${log.timestamp}-${index}`}
                className="flex items-center justify-between rounded-xl bg-zinc-900 p-4 font-mono text-sm"
              >
                <div>
                  <span className={colors[level] ?? "text-zinc-400"}>
                    [{level}]
                  </span>

                  <span className="ml-3 text-zinc-300">
                    {log.type}
                  </span>
                </div>

                <span className="text-xs text-zinc-500">
                  {new Date(log.timestamp).toLocaleTimeString()}
                </span>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}
