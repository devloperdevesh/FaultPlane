import { Cpu } from "lucide-react";
import MetricCard from "./MetricCard";

export default function WorkerCard() {
  return (
    <MetricCard
      title="Workers"
      value="12"
      description="Active compute workers"
      icon={Cpu}
      status="healthy"
      trend="+3 workers online"
    />
  );
}
