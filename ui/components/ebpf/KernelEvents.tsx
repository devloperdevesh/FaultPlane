"use client";

import { useEffect, useState } from "react";
import Panel from "@/components/ui/Panel";
import { getEbpfEvents, type EbpfEvent } from "@/lib/api";

export default function KernelEvents() {
  const [events, setEvents] = useState<EbpfEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    const load = async () => {
      try {
        const response = await getEbpfEvents();

        if (mounted) {
          setEvents(response.events);
          setError(null);
          setLoading(false);
        }
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load kernel events",
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
    <Panel>
      <div className="mb-4 flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-white">Kernel Events</h2>
          <p className="mt-1 text-xs text-zinc-500">
            Live events reported by the FaultPlane eBPF monitor
          </p>
        </div>

        <span className="rounded-full border border-white/10 px-3 py-1 text-[11px] text-zinc-500">
          {events.length} events
        </span>
      </div>

      {loading ? (
        <div className="py-8 text-center text-sm text-zinc-500">
          Loading kernel events...
        </div>
      ) : error ? (
        <div className="py-8 text-center">
          <p className="text-sm text-red-400">Kernel events unavailable</p>
          <p className="mt-2 text-xs text-zinc-600">{error}</p>
        </div>
      ) : events.length === 0 ? (
        <div className="py-8 text-center">
          <p className="text-sm text-zinc-400">No kernel events reported</p>
          <p className="mt-2 text-xs text-zinc-600">
            The eBPF monitor has not emitted any events yet.
          </p>
        </div>
      ) : (
        <div className="max-h-80 space-y-2 overflow-y-auto font-mono text-xs">
          {events.map((event, index) => (
            <div
              key={`${String(event.type ?? "event")}-${String(event.timestamp ?? index)}-${index}`}
              className="rounded-lg border border-white/[0.06] bg-zinc-900 p-3"
            >
              <div className="flex items-center justify-between gap-4">
                <span className="text-green-400">
                  {String(event.type ?? "UNKNOWN")}
                </span>

                {event.timestamp !== undefined && (
                  <span className="shrink-0 text-zinc-600">
                    {String(event.timestamp)}
                  </span>
                )}
              </div>

              {Object.keys(event).length > 2 && (
                <pre className="mt-2 overflow-x-auto whitespace-pre-wrap text-zinc-500">
                  {JSON.stringify(event, null, 2)}
                </pre>
              )}
            </div>
          ))}
        </div>
      )}
    </Panel>
  );
}
