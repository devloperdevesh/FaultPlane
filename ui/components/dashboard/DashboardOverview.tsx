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

import MemoryGrid from "@/components/memory/MemoryGrid";

import BlastGraph from "@/components/blast-radius/BlastGraph";
import FailurePropagation from "@/components/blast-radius/FailurePropagation";
import ImpactTimeline from "@/components/blast-radius/ImpactTimeline";
import StateDiff from "@/components/gitops/StateDiff";
import FinOpsOverview from "@/components/cards/FinOpsOverview";
import TenantTable from "@/components/tables/TenantTable";
import RuntimeShield from "@/components/status/RuntimeShield";
import CapabilityMatrix from "@/components/status/CapabilityMatrix";
import ThreatScore from "@/components/status/ThreatScore";
import InterceptorLogs from "@/components/telemetry/InterceptorLogs";
import ProxyStream from "@/components/telemetry/ProxyStream";
import HijackDetection from "@/components/telemetry/HijackDetection";
import TransportMetrics from "@/components/telemetry/TransportMetrics";
import SpeculativePath from "@/components/workflows/SpeculativePath";

function Section({
  title,
  description,
  children,
}: {
  title: string;
  description?: string;
  children: React.ReactNode;
}) {
  return (
    <section
      className="
      rounded-2xl
      border
      border-white/10
      bg-zinc-950/60
      backdrop-blur-xl
      p-6
      shadow-xl
      shadow-black/20
      "
    >
      <div className="mb-6">
        <h2
          className="
          text-sm
          font-semibold
          text-white
          "
        >
          {title}
        </h2>

        {description && (
          <p
            className="
            mt-1
            text-xs
            text-zinc-500
            "
          >
            {description}
          </p>
        )}
      </div>

      {children}
    </section>
  );
}

export default function DashboardPage() {
  return (
    <div
      className="
space-y-8
"
    >
      {/* Header */}

      <header>
        <h1
          className="
text-3xl
font-semibold
tracking-tight
text-white
"
        >
          Operations Console
        </h1>

        <p
          className="
mt-2
text-sm
text-zinc-500
"
        >
          FaultPlane runtime observability, recovery intelligence and transport
          resilience control plane.
        </p>
      </header>

      {/* Runtime Overview */}

      <Section
        title="Runtime Overview"
        description="Current infrastructure health and recovery metrics"
      >
        <div
          className="
grid
gap-4
md:grid-cols-2
xl:grid-cols-5
"
        >
          <MetricCard title="Gateway" value="Healthy" />

          <MetricCard title="Active Workers" value="24" />

          <RecoveryCard />

          <CostCard />

          <SLAProtectionCard />
        </div>
      </Section>

      {/* Performance */}

      <Section
        title="Performance Telemetry"
        description="Latency boundaries and recovery execution metrics"
      >
        <div
          className="
grid
gap-6
xl:grid-cols-2
"
        >
          <LatencyChart />

          <RecoveryTimeline />
        </div>
      </Section>

      {/* Workload */}

      <Section
        title="Workload Runtime"
        description="Agent execution and worker state management"
      >
        <div
          className="
space-y-6
"
        >
          <WorkersTable />

          <WorkflowTable />
        </div>
      </Section>

      {/* Observability */}

      <Section
        title="Kernel Observability"
        description="Low level runtime events and system telemetry"
      >
        <div
          className="
grid
gap-6
xl:grid-cols-2
"
        >
          <TelemetryLogs />

          <KernelLogs />
        </div>

        <div
          className="
mt-6
"
        >
          <EbpfMonitor />
        </div>
      </Section>

      {/* Infrastructure */}

      <Section
        title="Transport Infrastructure"
        description="Network failover and connection migration state"
      >
        <SocketMigration />
      </Section>

      {/* Memory */}

      <Section
        title="State Resilience"
        description="Runtime memory migration and checkpoint restoration"
      >
        <MemoryGrid />
      </Section>

      {/* Recovery Intelligence */}

      <Section
        title="Recovery Intelligence"
        description="Failure impact analysis and blast radius containment"
      >
        <div
          className="
grid
gap-6
xl:grid-cols-3
"
        >
          <BlastGraph />

          <FailurePropagation />

          <ImpactTimeline />
        </div>

        <div
          className="
mt-6
"
        >
          <CheckpointTimeline />
        </div>
      </Section>
      {/* GitOps Debugging */}

      <Section
        title="Checkpoint Debugging"
        description="Runtime state comparison and recovery verification"
      >
        <StateDiff />
      </Section>
      {/* FinOps Intelligence */}

      <Section
        title="FinOps Intelligence"
        description="Cloud cost optimization and resource efficiency"
      >
        <FinOpsOverview />
      </Section>

      {/* Multi Tenant Governance */}

      <Section
        title="Multi Tenant Governance"
        description="Tenant isolation, quota and resource management"
      >
        <TenantTable />
      </Section>

      <Section
        title="Runtime Security"
        description="WASM sandbox isolation and workload protection"
      >
        <div
          className="
space-y-6
"
        >
          <RuntimeShield />

          <div
            className="
grid
gap-6
xl:grid-cols-2
"
          >
            <CapabilityMatrix />

            <ThreatScore />
          </div>
        </div>
      </Section>

      <Section
        title="Network Security Telemetry"
        description="Transport interception and connection intelligence"
      >
        <div className="space-y-6">
          <TransportMetrics />

          <div
            className="
grid
gap-6
xl:grid-cols-2
"
          >
            <InterceptorLogs />

            <ProxyStream />
          </div>

          <HijackDetection />
        </div>
      </Section>
      {/* Speculative Execution */}

      <Section
        title="Workflow Intelligence"
        description="Speculative execution and branch prediction"
      >
        <SpeculativePath />
      </Section>
    </div>
  );
}
