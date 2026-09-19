"use client";

import { useEffect, useState } from "react";
import Panel from "@/components/ui/Panel";
import { getEbpfHooks, type EbpfHook } from "@/lib/api";

export default function HookStatus() {
  const [hooks, setHooks] = useState<EbpfHook[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    const load = async () => {
      try {
        const response = await getEbpfHooks();

        if (mounted) {
          setHooks(response.hooks);
          setError(null);
          setLoading(false);
        }
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load eBPF hooks",
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
      <h2 className="mb-4 font-semibold text-white">eBPF Kernel Hooks</h2>

      {loading ? (
        <div className="py-6 text-center text-sm text-zinc-500">
          Loading kernel hooks...
        </div>
      ) : error ? (
        <div className="py-6 text-center">
          <p className="text-sm text-red-400">Kernel hooks unavailable</p>
          <p className="mt-2 text-xs text-zinc-600">{error}</p>
        </div>
      ) : hooks.length === 0 ? (
        <div className="py-6 text-center">
          <p className="text-sm text-zinc-400">No eBPF hooks reported</p>
          <p className="mt-2 text-xs text-zinc-600">
            FaultPlane has not reported any active kernel hooks yet.
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {hooks.map((hook, index) => (
            <div
              key={`${hook.name}-${index}`}
              className="flex justify-between rounded-lg bg-zinc-900 p-3"
            >
              <span className="text-zinc-300">{hook.name}</span>

              <span className="font-mono text-sm text-green-400">
                {hook.status ?? "UNKNOWN"}
              </span>
            </div>
          ))}
        </div>
      )}
    </Panel>
  );
}
