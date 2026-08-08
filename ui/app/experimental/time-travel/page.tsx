import RoutePage from "@/components/layout/RoutePage";
import StateDiffViewer from "@/components/experimental/time-travel/StateDiffViewer";
import SnapshotTimeline from "@/components/experimental/time-travel/SnapshotTimeline";

export default function TimeTravelPage() {
  return (
    <RoutePage
      eyebrow="Experiment"
      title="Time Travel"
      description="Inspect snapshot deltas and historical state transitions in the experimental time-travel workspace."
    >
      <div className="space-y-6">
        <SnapshotTimeline />
        <StateDiffViewer />
      </div>
    </RoutePage>
  );
}
