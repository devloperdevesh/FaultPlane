import {
  Activity,
  ShieldCheck,
  Workflow,
  Database,
  Cpu,
  Clock3,
} from "lucide-react";

export default function Home() {
  return (
    <main className="min-h-screen bg-zinc-950 text-white p-8">
      <div className="mx-auto max-w-7xl">
        <div className="mb-10">
          <h1 className="text-4xl font-bold">
            FaultPlane Operations Dashboard
          </h1>

          <p className="mt-2 text-zinc-400">
            Control Plane for AI Agent Reliability
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          <MetricCard
            title="Gateway Status"
            value="Healthy"
            icon={<ShieldCheck size={28} />}
          />

          <MetricCard
            title="Workers"
            value="12 Active"
            icon={<Cpu size={28} />}
          />

          <MetricCard
            title="Workflows"
            value="37 Running"
            icon={<Workflow size={28} />}
          />

          <MetricCard
            title="Checkpoints"
            value="126 Saved"
            icon={<Database size={28} />}
          />

          <MetricCard
            title="Recovery Events"
            value="18"
            icon={<Activity size={28} />}
          />

          <MetricCard
            title="Latency"
            value="1.8 ms"
            icon={<Clock3 size={28} />}
          />
        </div>
      </div>
    </main>
  );
}

type MetricCardProps = {
  title: string;
  value: string;
  icon: React.ReactNode;
};

function MetricCard({ title, value, icon }: MetricCardProps) {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-900 p-6">
      <div className="mb-4 text-blue-400">{icon}</div>

      <h2 className="text-sm uppercase tracking-wide text-zinc-400">{title}</h2>

      <p className="mt-2 text-3xl font-bold">{value}</p>
    </div>
  );
}
