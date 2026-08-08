import RoutePage from "@/components/layout/RoutePage";
import WorkerCard from "@/components/cards/WorkerCard";
import LatencyCard from "@/components/cards/LatencyCard";
import WorkerTable from "@/components/workers/WorkerTable";

export default function WorkersPage() {
  return (
    <RoutePage
      eyebrow="Fleet"
      title="Workers"
      description="Monitor the active compute fleet, worker health, and execution distribution across the runtime layer."
    >
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <WorkerCard />
          <LatencyCard />
        </div>

        <WorkerTable />
      </div>
    </RoutePage>
  );
}
