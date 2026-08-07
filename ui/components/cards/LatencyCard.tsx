import { Activity } from "lucide-react";
import MetricCard from "./MetricCard";

export default function LatencyCard() {
  return (
    <MetricCard
      title="Latency P99"
      value="2.3ms"
      description="Runtime response latency"
      icon={Activity}
      status="healthy"
      trend="↓ 18% improved"
    />
  );
}
