import RoutePage from "@/components/layout/RoutePage";
import TransportMetrics from "@/components/telemetry/TransportMetrics";
import ProxyStream from "@/components/telemetry/ProxyStream";
import InterceptorLogs from "@/components/telemetry/InterceptorLogs";

export default function TelemetryPage() {
  return (
    <RoutePage
      eyebrow="Signals"
      title="Telemetry"
      description="Observe transport metrics, proxy streams, and interception logs from the runtime instrumentation layer."
    >
      <div className="space-y-6">
        <TransportMetrics />

        <div className="grid gap-6 xl:grid-cols-2">
          <ProxyStream />
          <InterceptorLogs />
        </div>
      </div>
    </RoutePage>
  );
}
