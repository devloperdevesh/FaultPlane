"use client";

import GatewayCard from "@/components/cards/GatewayCard";
import WorkerCard from "@/components/cards/WorkerCard";
import RecoveryCard from "@/components/cards/RecoveryCard";
import LatencyCard from "@/components/cards/LatencyCard";
import CheckpointCard from "@/components/cards/CheckpointCard";

import LatencyChart from "@/components/charts/LatencyChart";
import RecoveryTimeline from "@/components/charts/RecoveryTimeline";

import WorkersTable from "@/components/tables/WorkersTable";
import WorkflowTable from "@/components/tables/WorkflowTable";

import TelemetryLogs from "@/components/logs/TelemetryLogs";
import KernelLogs from "@/components/logs/KernelLogs";

import EbpfMonitor from "@/components/ebpf/EbpfMonitor";

import SocketMigration from "@/components/topology/SocketMigration";

import MemoryGrid from "@/components/memory/MemoryGrid";

import BlastGraph from "@/components/blast-radius/BlastGraph";
import FailurePropagation from "@/components/blast-radius/FailurePropagation";
import ImpactTimeline from "@/components/blast-radius/ImpactTimeline";

import CheckpointTimeline from "@/components/gitops/CheckpointTimeline";
import StateDiff from "@/components/gitops/StateDiff";

import FinOpsOverview from "@/components/cards/FinOpsOverview";

import TenantTable from "@/components/tables/TenantTable";

import RuntimeShield from "@/components/status/RuntimeShield";
import CapabilityMatrix from "@/components/status/CapabilityMatrix";
import ThreatScore from "@/components/status/ThreatScore";

import TransportMetrics from "@/components/telemetry/TransportMetrics";
import InterceptorLogs from "@/components/telemetry/InterceptorLogs";
import ProxyStream from "@/components/telemetry/ProxyStream";
import HijackDetection from "@/components/telemetry/HijackDetection";

import SpeculativePath from "@/components/workflows/SpeculativePath";

function Section({
  title,
  description,
  children,
}: {
  title: string;
  description: string;
  children: React.ReactNode;
}) {
  return (
    <section
      className="
rounded-2xl
border
border-white/10
bg-zinc-950/70
backdrop-blur-xl
p-6
shadow-xl
shadow-black/20
"
    >
      <div
        className="
mb-6
"
      >
        <h2
          className="
text-sm
font-semibold
tracking-wide
text-white
"
        >
          {title}
        </h2>

        <p
          className="
mt-1
text-xs
text-zinc-500
"
        >
          {description}
        </p>
      </div>

      {children}
    </section>
  );
}

export default function DashboardOverview() {
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
        description="Live infrastructure health and runtime metrics"
      >
        <div
          className="
grid
gap-4
md:grid-cols-2
xl:grid-cols-5
"
        >
          <GatewayCard />

          <WorkerCard />

          <RecoveryCard />

          <LatencyCard />

          <CheckpointCard />
        </div>
      </Section>

      {/* Performance */}

      <Section
        title="Performance Telemetry"
        description="Latency boundaries and recovery execution"
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

      {/* Workloads */}

      <Section
        title="Workload Runtime"
        description="Worker execution and workflow orchestration"
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

      {/* Kernel */}

      <Section
        title="Kernel Observability"
        description="Low level runtime signals and system events"
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

        <div className="mt-6">
          <EbpfMonitor />
        </div>
      </Section>

      {/* Infrastructure */}

      <Section
        title="Transport Infrastructure"
        description="Network migration and connection resilience"
      >
        <SocketMigration />
      </Section>

      {/* Memory */}

      <Section
        title="State Resilience"
        description="Checkpoint recovery and memory migration"
      >
        <MemoryGrid />
      </Section>

      {/* Recovery Intelligence */}

      <Section
        title="Recovery Intelligence"
        description="Failure analysis and blast radius containment"
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

      {/* Debugging */}

      <Section
        title="Checkpoint Debugging"
        description="Runtime state comparison and recovery verification"
      >
        <StateDiff />
      </Section>

      {/* FinOps */}

      <Section
        title="FinOps Intelligence"
        description="Cloud cost optimization and efficiency"
      >
        <FinOpsOverview />
      </Section>

      {/* Multi Tenant */}

      <Section
        title="Multi Tenant Governance"
        description="Tenant isolation and resource control"
      >
        <TenantTable />
      </Section>

      {/* Security */}

      <Section
        title="Runtime Security"
        description="Sandbox isolation and workload protection"
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

      {/* Network Security */}

      <Section
        title="Network Security Telemetry"
        description="Transport interception intelligence"
      >
        <div
          className="
space-y-6
"
        >
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

      {/* Workflow Intelligence */}

      <Section
        title="Workflow Intelligence"
        description="Speculative execution and branch prediction"
      >
        <SpeculativePath />
      </Section>
    </div>
  );
}
