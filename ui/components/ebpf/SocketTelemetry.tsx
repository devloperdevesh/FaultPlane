"use client";

import { useEffect, useState } from "react";
import Panel from "@/components/ui/Panel";
import { getTelemetry, type TelemetryEvent } from "@/lib/api";

export default function SocketTelemetry() {
  const [events, setEvents] = useState<TelemetryEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    const load = async () => {
      try {
        const response = await getTelemetry();

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
              : "Unable to load socket telemetry",
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

  const httpEvents = events.filter(
    (event) =>
      event.type === "http_request" ||
      event.type === "request_complete",
  );

  const networkEvents = events.filter(
    (event) =>
      event.type.toLowerCase().includes("network") ||
      event.type.toLowerCase().includes("socket") ||
      event.type.toLowerCase().includes("tcp") ||
      event.type.toLowerCase().includes("connection"),
  );

  return (
    <Panel>
      <div className="mb-4">
        <h2 className="font-semibold text-white">Socket Telemetry</h2>
        <p className="mt-1 text-xs text-zinc-500">
          Live network-related telemetry reported by FaultPlane
        </p>
      </div>

      {loading ? (
        <div className="py-8 text-center text-sm text-zinc-500">
          Loading socket telemetry...
        </div>
      ) : error ? (
        <div className="py-8 text-center">
          <p className="text-sm text-red-400">
            Socket telemetry unavailable
          </p>
          <p className="mt-2 text-xs text-zinc-600">{error}</p>
        </div>
      ) : events.length === 0 ? (
        <div className="py-8 text-center">
          <p className="text-sm text-zinc-400">
            No socket telemetry reported
          </p>
          <p className="mt-2 text-xs text-zinc-600">
            FaultPlane has not reported network telemetry yet.
          </p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 gap-4">
            <div className="rounded-lg bg-zinc-900 p-4">
              <p className="text-xs text-zinc-500">Telemetry Events</p>
              <p className="mt-2 text-xl text-white">{events.length}</p>
            </div>

            <div className="rounded-lg bg-zinc-900 p-4">
              <p className="text-xs text-zinc-500">HTTP Events</p>
              <p className="mt-2 text-xl text-white">{httpEvents.length}</p>
            </div>

            <div className="rounded-lg bg-zinc-900 p-4">
              <p className="text-xs text-zinc-500">Network Events</p>
              <p className="mt-2 text-xl text-white">{networkEvents.length}</p>
            </div>

            <div className="rounded-lg bg-zinc-900 p-4">
              <p className="text-xs text-zinc-500">Latest Event</p>
              <p className="mt-2 truncate font-mono text-sm text-white">
                {events.at(-1)?.type ?? "—"}
              </p>
            </div>
          </div>

          <div className="mt-4 max-h-64 space-y-2 overflow-y-auto">
            {events.slice(-20).reverse().map((event, index) => (
              <div
                key={`${event.type}-${event.timestamp}-${index}`}
                className="rounded-lg border border-white/[0.06] bg-zinc-900 p-3"
              >
                <div className="flex items-center justify-between gap-4">
                  <span className="font-mono text-xs text-zinc-300">
                    {event.type}
                  </span>

                  <span className="shrink-0 text-[11px] text-zinc-600">
                    {new Date(event.timestamp).toLocaleTimeString()}
                  </span>
                </div>

                {event.value !== undefined && (
                  <p className="mt-1 text-xs text-zinc-500">
                    value: {event.value}
                  </p>
                )}
              </div>
            ))}
          </div>
        </>
      )}
    </Panel>
  );
}
