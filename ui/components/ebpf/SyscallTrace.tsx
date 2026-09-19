"use client";

import { useEffect, useState } from "react";
import { getEbpfEvents, getEbpfHooks } from "@/lib/api";
import type { EbpfEvent, EbpfHook } from "@/lib/types";

export default function SyscallTrace() {
  const [hooks, setHooks] = useState<EbpfHook[]>([]);
  const [events, setEvents] = useState<EbpfEvent[]>([]);
  const [error, setError] = useState(false);

  useEffect(() => {
    let mounted = true;

    async function loadData() {
      try {
        const [hookData, eventData] = await Promise.all([
          getEbpfHooks(),
          getEbpfEvents(),
        ]);

        if (mounted) {
          setHooks(hookData.hooks);
          setEvents(eventData.events);
          setError(false);
        }
      } catch {
        if (mounted) {
          setError(true);
        }
      }
    }

    void loadData();

    const interval = window.setInterval(loadData, 2000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  return (
    <section className="rounded-xl border border-white/10 bg-zinc-950 p-6">
      <div className="mb-6">
        <h2 className="text-sm font-semibold text-white">
          eBPF Runtime Trace
        </h2>
        <p className="mt-1 text-xs text-zinc-500">
          Live hook and kernel event visibility
        </p>
      </div>

      {error ? (
        <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
          eBPF telemetry unavailable.
        </div>
      ) : hooks.length === 0 && events.length === 0 ? (
        <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
          No eBPF hooks or events reported.
        </div>
      ) : (
        <div className="space-y-6">
          <div>
            <p className="mb-3 text-xs uppercase tracking-wider text-zinc-500">
              Hooks
            </p>

            <div className="space-y-2">
              {hooks.map((hook, index) => (
                <div
                  key={`${hook.name}-${index}`}
                  className="flex items-center justify-between rounded-lg border border-white/10 bg-zinc-900 p-3"
                >
                  <span className="font-mono text-xs text-zinc-300">
                    {hook.name}
                  </span>

                  <span className="text-xs text-zinc-400">
                    {hook.status ?? "reported"}
                  </span>
                </div>
              ))}
            </div>
          </div>

          <div>
            <p className="mb-3 text-xs uppercase tracking-wider text-zinc-500">
              Events
            </p>

            <div className="space-y-2">
              {events.slice(0, 10).map((event, index) => (
                <pre
                  key={index}
                  className="overflow-x-auto rounded-lg border border-white/10 bg-black p-3 text-xs text-zinc-400"
                >
                  {JSON.stringify(event, null, 2)}
                </pre>
              ))}
            </div>
          </div>
        </div>
      )}
    </section>
  );
}
