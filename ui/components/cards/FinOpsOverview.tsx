import CostSavings from "./CostSavings";
import SpotComparison from "./SpotComparison";
import OptimizationScore from "./OptimizationScore";
import ResourceWaste from "./ResourceWaste";

export default function FinOpsOverview() {
  return (
    <div
      className="
space-y-6
"
    >
      <CostSavings />

      <div
        className="
grid
gap-6
xl:grid-cols-3
"
      >
        <SpotComparison />

        <OptimizationScore />

        <ResourceWaste />
      </div>
    </div>
  );
}
