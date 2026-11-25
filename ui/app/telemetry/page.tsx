import RoutePage from "@/components/layout/RoutePage";
import TransportMetrics from "@/components/telemetry/TransportMetrics";
import ProxyStream from "@/components/telemetry/ProxyStream";
import InterceptorLogs from "@/components/telemetry/InterceptorLogs";
import ComplianceTelemetryPanel from "@/components/telemetry/ComplianceTelemetryPanel";

export default function TelemetryPage() {
  return (
    <RoutePage
      eyebrow="Signals"
      title="Telemetry & Compliance"
      description="Observe transport metrics, proxy streams, verifiable compliance telemetry, and export cryptographically signed audit ledgers."
    >
      <div className="space-y-8">
        <ComplianceTelemetryPanel />

        <TransportMetrics />

        <div className="grid gap-6 xl:grid-cols-2">
          <ProxyStream />
          <InterceptorLogs />
        </div>
      </div>
    </RoutePage>
  );
}
