import { Server } from "lucide-react";
import MetricCard from "./MetricCard";

export default function GatewayCard() {
  return (
    <MetricCard
      title="Gateway"
      value="Healthy"
      description="Ingress runtime status"
      icon={Server}
      status="healthy"
      trend="99.99% uptime"
    />
  );
}
