import RoutePage from "@/components/layout/RoutePage";
import CheckpointTable from "@/components/tables/CheckpointTable";
import CheckpointCard from "@/components/cards/CheckpointCard";
import CheckpointTimeline from "@/components/gitops/CheckpointTimeline";

export default function CheckpointsPage() {
  return (
    <RoutePage
      eyebrow="State"
      title="Checkpoints"
      description="Inspect recovery snapshots, checkpoint cadence, and recent state persistence activity."
    >
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <CheckpointCard />
        </div>

        <CheckpointTimeline />

        <CheckpointTable />
      </div>
    </RoutePage>
  );
}
