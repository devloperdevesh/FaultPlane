"use client";

import { useEffect, useState } from "react";
import { getNetworkEvents, type TelemetryEvent } from "@/lib/api";

export default function KernelLogs() {
  const [events, setEvents] = useState<TelemetryEvent[]>([]);
  const [error, setError] = useState(false);

  useEffect(() => {
    let active = true;

    const load = async () => {
      try {
        const response = await getNetworkEvents();
        if (active) {
          setEvents(response.events);
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
      <div className="mb-5 flex items-center justify-between">
        <h2 className="text-sm font-semibold text-white">
          Kernel Events
        </h2>

        <span className="rounded-full bg-blue-500/10 px-3 py-1 text-xs text-blue-400">
          {error ? "UNAVAILABLE" : "LIVE"}
        </span>
      </div>

      <div className="space-y-3">
        {events.length === 0 ? (
          <div className="rounded-xl bg-zinc-900 p-4 text-sm text-zinc-500">
            {error
              ? "Network events API unavailable."
              : "No kernel events available."}
          </div>
        ) : (
          events.map((event, index) => (
            <div
              key={`${event.timestamp}-${index}`}
              className="flex items-center justify-between rounded-xl bg-zinc-900 p-4 font-mono text-xs"
            >
              <div>
                <span className="mr-3 text-blue-400">{">"}</span>
                <span className="text-zinc-300">
                  {event.type}
                </span>
              </div>

              <span className="text-zinc-500">
                {new Date(event.timestamp).toLocaleTimeString()}
              </span>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
