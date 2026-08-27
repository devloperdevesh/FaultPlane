"use client";

import { useEffect, useState } from "react";
import { Activity, Cpu, Database, Server } from "lucide-react";
import MetricCard from "../cards/MetricCard";
import { getMetrics } from "../../lib/api";
import type { DashboardMetrics } from "../../lib/types";

export default function DashboardOverview() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    const loadMetrics = async () => {
      try {
        const next = await getMetrics();

        if (mounted) {
          setMetrics(next);
          setError(null);
        }
      } catch (err) {
        if (mounted) {
          setError(
            err instanceof Error
              ? err.message
              : "Unable to load FaultPlane runtime metrics",
          );
        }
      }
    };

    void loadMetrics();

    const interval = window.setInterval(() => {
      void loadMetrics();
    }, 2000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  if (error) {
    return (
      <section className="rounded-xl border border-red-500/20 bg-zinc-950 p-6">
        <p className="text-sm font-medium text-red-400">
          Runtime metrics unavailable
        </p>
        <p className="mt-2 text-xs text-zinc-500">{error}</p>
      </section>
    );
  }

  if (!metrics) {
    return (
      <section className="rounded-xl border border-white/10 bg-zinc-950 p-6">
        <p className="text-sm text-zinc-400">
          Loading live FaultPlane runtime metrics...
        </p>
      </section>
    );
  }

  return (
    <section className="space-y-6">
      <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard
          title="Requests"
          value={metrics.requests.toLocaleString()}
          description="Runtime requests recorded"
          icon={Activity}
          status="healthy"
        />

        <MetricCard
          title="Average Latency"
          value={`${metrics.latency.toFixed(2)} ms`}
          description="Runtime average latency"
          icon={Activity}
          status="healthy"
        />

        <MetricCard
          title="CPU"
          value={
            metrics.cpu > 0
              ? `${metrics.cpu.toFixed(2)}%`
              : "Unavailable"
          }
          description="Runtime CPU metric"
          icon={Cpu}
          status={metrics.cpu > 0 ? "healthy" : "warning"}
        />

        <MetricCard
          title="Memory"
          value={`${metrics.memory.toFixed(2)} MB`}
          description="Runtime memory usage"
          icon={Database}
          status="healthy"
        />
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <MetricCard
          title="Workers"
          value={metrics.workers.toLocaleString()}
          description="Workers reported by runtime telemetry"
          icon={Server}
          status={metrics.workers > 0 ? "healthy" : "warning"}
        />

        <MetricCard
          title="Recoveries"
          value={metrics.recoveries.toLocaleString()}
          description="Recorded runtime recovery events"
          icon={Activity}
          status="healthy"
        />
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-950 p-5">
        <p className="text-xs uppercase tracking-widest text-zinc-500">
          Runtime snapshot
        </p>

        <div className="mt-4 grid gap-4 sm:grid-cols-3">
          <div>
            <p className="text-xs text-zinc-500">Checkpoints</p>
            <p className="mt-1 text-xl font-semibold text-white">
              {metrics.checkpoints.toLocaleString()}
            </p>
          </div>

          <div>
            <p className="text-xs text-zinc-500">Updated</p>
            <p className="mt-1 text-sm text-zinc-300">
              {new Date(metrics.updatedAt).toLocaleString()}
            </p>
          </div>

          <div>
            <p className="text-xs text-zinc-500">Source</p>
            <p className="mt-1 text-sm text-emerald-400">
              FaultPlane runtime API
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
