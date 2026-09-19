"use client";

import { useEffect, useState } from "react";
import { getEbpfHooks, getEbpfStatus } from "@/lib/api";
import type { EbpfHook } from "@/lib/types";

export default function KernelMonitor() {
  const [status, setStatus] = useState<string | null>(null);
  const [hooks, setHooks] = useState<EbpfHook[]>([]);
  const [error, setError] = useState(false);

  useEffect(() => {
    let mounted = true;

    async function load() {
      try {
        const [statusData, hookData] = await Promise.all([
          getEbpfStatus(),
          getEbpfHooks(),
        ]);

        if (mounted) {
          setStatus(statusData.status);
          setHooks(hookData.hooks);
          setError(false);
        }
      } catch {
        if (mounted) {
          setError(true);
        }
      }
    }

    void load();

    const interval = window.setInterval(load, 2000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  return (
    <div className="rounded-xl border border-zinc-800 bg-black p-5">
      <h2 className="mb-5 text-sm font-semibold text-white">
        Kernel Monitor
      </h2>

      {error ? (
        <div className="rounded-lg border border-white/10 bg-zinc-900 p-3 text-sm text-zinc-500">
          Kernel telemetry unavailable.
        </div>
      ) : status || hooks.length > 0 ? (
        <div className="space-y-3">
          <div className="flex justify-between rounded-lg bg-zinc-900 px-3 py-2">
            <span className="font-mono text-xs text-zinc-400">
              eBPF Status
            </span>
            <span className="font-mono text-xs text-zinc-300">
              {status ?? "—"}
            </span>
          </div>

          <div className="flex justify-between rounded-lg bg-zinc-900 px-3 py-2">
            <span className="font-mono text-xs text-zinc-400">
              Reported Hooks
            </span>
            <span className="font-mono text-xs text-zinc-300">
              {hooks.length}
            </span>
          </div>
        </div>
      ) : (
        <div className="rounded-lg border border-white/10 bg-zinc-900 p-3 text-sm text-zinc-500">
          No kernel telemetry reported.
        </div>
      )}
    </div>
  );
}
