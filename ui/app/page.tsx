import React from "react";
import {
  Activity,
  ShieldCheck,
  Workflow,
  Database,
  Cpu,
  Clock3,
} from "lucide-react";

import AppShell from "@/components/layout/AppShell";

type MetricCardProps = {
  title: string;
  value: string;
  icon: React.ReactNode;
};

export default function Home() {
  return (
    <AppShell>
      <div className="mx-auto max-w-7xl p-8">
        {/* Header */}
        <header className="mb-10">
          <h1 className="text-4xl font-bold tracking-tight text-white">
            FaultPlane Operations Dashboard
          </h1>

          <p className="mt-2 text-zinc-400">
            Control Plane for AI Agent Reliability
          </p>
        </header>

        {/* Metrics */}
        <section className="grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
          <MetricCard
            title="Gateway Status"
            value="Healthy"
            icon={<ShieldCheck size={26} />}
          />

          <MetricCard
            title="Workers"
            value="12 Active"
            icon={<Cpu size={26} />}
          />

          <MetricCard
            title="Workflows"
            value="37 Running"
            icon={<Workflow size={26} />}
          />

          <MetricCard
            title="Checkpoints"
            value="126 Saved"
            icon={<Database size={26} />}
          />

          <MetricCard
            title="Recovery Events"
            value="18"
            icon={<Activity size={26} />}
          />

          <MetricCard
            title="Latency"
            value="1.8 ms"
            icon={<Clock3 size={26} />}
          />
        </section>
      </div>
    </AppShell>
  );
}

function MetricCard({ title, value, icon }: MetricCardProps) {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-900/70 p-6 transition-all duration-300 hover:border-blue-500 hover:bg-zinc-900 hover:shadow-lg hover:shadow-blue-500/10">
      <div className="mb-5 flex h-12 w-12 items-center justify-center rounded-lg bg-blue-500/10 text-blue-400">
        {icon}
      </div>

      <p className="text-xs font-medium uppercase tracking-widest text-zinc-500">
        {title}
      </p>

      <h2 className="mt-3 text-3xl font-bold text-white">{value}</h2>
    </div>
  );
}
