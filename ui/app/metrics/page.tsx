import RoutePage from "@/components/layout/RoutePage";
import CpuChart from "@/components/charts/CpuChart";
import MemoryChart from "@/components/charts/MemoryChart";
import LatencyChart from "@/components/charts/LatencyChart";

export default function MetricsPage() {
  return (
    <RoutePage
      eyebrow="Performance"
      title="Metrics"
      description="Review core runtime indicators for CPU pressure, memory usage, and latency distribution."
    >
      <div className="grid gap-6 xl:grid-cols-2">
        <CpuChart />
        <MemoryChart />
        <LatencyChart />
      </div>
    </RoutePage>
  );
}
