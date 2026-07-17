import MetricCard from "@/components/cards/MetricCard";
import SLAProtectionCard from "@/components/cards/SLAProtectionCard";
import RecoveryCard from "@/components/cards/RecoveryCard";
import CostCard from "@/components/cards/CostCard";

import LatencyChart from "@/components/charts/LatencyChart";
import RecoveryTimeline from "@/components/charts/RecoveryTimeline";

import WorkersTable from "@/components/tables/WorkersTable";
import WorkflowTable from "@/components/tables/WorkflowTable";

import TelemetryLogs from "@/components/logs/TelemetryLogs";
import KernelLogs from "@/components/logs/KernelLogs";

import EbpfMonitor from "@/components/ebpf/EbpfMonitor";
import SocketMigration from "@/components/topology/SocketMigration";
import CheckpointTimeline from "@/components/gitops/CheckpointTimeline";

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      {/* Header */}

      <section>
        <h1
          className="
          text-2xl
          font-semibold
          text-white
          "
        >
          Operations Dashboard
        </h1>

        <p
          className="
          mt-1
          text-sm
          text-zinc-500
          "
        >
          FaultPlane runtime observability and recovery control
        </p>
      </section>

      {/* Runtime Metrics */}

      <section
        className="
        grid
        gap-4
        md:grid-cols-2
        xl:grid-cols-5
        "
      >
        <MetricCard title="Gateway" value="Healthy" />

        <MetricCard title="Active Workers" value="2" />

        <RecoveryCard />

        <CostCard />

        <SLAProtectionCard />
      </section>

      {/* Runtime Charts */}

      <section
        className="
        grid
        gap-6
        xl:grid-cols-2
        "
      >
        <LatencyChart />

        <RecoveryTimeline />
      </section>

      {/* Workload Tables */}

      <section
        className="
        space-y-6
        "
      >
        <WorkersTable />

        <WorkflowTable />
      </section>

      {/* Logs */}

      <section
        className="
        grid
        gap-6
        xl:grid-cols-2
        "
      >
        <TelemetryLogs />

        <KernelLogs />
      </section>

      {/* Recovery History */}

      <section>
        <CheckpointTimeline />
      </section>

      {/* Kernel Layer */}

      <section>
        <EbpfMonitor />
      </section>

      {/* Transport Layer */}

      <section>
        <SocketMigration />
      </section>
    </div>
  );
}
