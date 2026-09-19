"use client";

import { useEffect, useState } from "react";
import { getHealth } from "@/lib/api";
import type { HealthResponse } from "@/lib/types";

export default function GatewayCard() {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    let mounted = true;

    async function loadHealth() {
      try {
        const data = await getHealth();

        if (mounted) {
          setHealth(data);
          setError(false);
        }
      } catch {
        if (mounted) {
          setError(true);
        }
      }
    }

    void loadHealth();

    const interval = window.setInterval(loadHealth, 2000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  const status = health?.status ?? (error ? "Unavailable" : "—");

  return (
    <section className="rounded-xl border border-white/10 bg-zinc-950 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-sm font-semibold text-white">
            Gateway
          </h2>
          <p className="mt-1 text-xs text-zinc-500">
            Runtime gateway health
          </p>
        </div>

        <span className="text-sm text-zinc-300">
          {status}
        </span>
      </div>

      <div className="mt-6 rounded-lg border border-white/10 bg-zinc-900/60 p-4">
        <p className="text-xs text-zinc-500">
          Uptime
        </p>
        <p className="mt-2 font-mono text-lg text-white">
          —
        </p>
      </div>
    </section>
  );
}
