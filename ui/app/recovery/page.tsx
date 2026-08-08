import RoutePage from "@/components/layout/RoutePage";
import RecoveryCard from "@/components/cards/RecoveryCard";
import RecoveryTimeline from "@/components/charts/RecoveryTimeline";
import BlastRadius from "@/components/blast-radius/BlastRadius";

export default function RecoveryPage() {
  return (
    <RoutePage
      eyebrow="Resilience"
      title="Recovery"
      description="Follow failover activity, recovery latency, and blast radius containment as incidents unfold."
    >
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          <RecoveryCard />
        </div>

        <RecoveryTimeline />

        <BlastRadius />
      </div>
    </RoutePage>
  );
}
