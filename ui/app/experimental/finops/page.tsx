import RoutePage from "@/components/layout/RoutePage";
import CostArbitrage from "@/components/experimental/finops/CostArbitrage";
import FinOpsOverview from "@/components/cards/FinOpsOverview";

export default function FinOpsPage() {
  return (
    <RoutePage
      eyebrow="Experiment"
      title="FinOps"
      description="Review compute economics, arbitrage opportunities, and cost-saving signals for the runtime layer."
    >
      <div className="space-y-6">
        <FinOpsOverview />
        <CostArbitrage />
      </div>
    </RoutePage>
  );
}
