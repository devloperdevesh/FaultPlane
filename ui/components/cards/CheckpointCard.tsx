import { Database } from "lucide-react";
import MetricCard from "./MetricCard";

export default function CheckpointCard() {
  return (
    <MetricCard
      title="Checkpoints"
      value="248"
      description="Recovery snapshots"
      icon={Database}
      status="healthy"
      trend="Last sync 2s ago"
    />
  );
}
