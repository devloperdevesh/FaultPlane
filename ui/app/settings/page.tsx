import RoutePage from "@/components/layout/RoutePage";
import RuntimeShield from "@/components/status/RuntimeShield";
import SystemResources from "@/components/status/SystemResources";
import CapabilityMatrix from "@/components/status/CapabilityMatrix";

export default function SettingsPage() {
  return (
    <RoutePage
      eyebrow="Control"
      title="Settings"
      description="Tune runtime posture, feature coverage, and system guardrails from a single place."
    >
      <div className="space-y-6">
        <RuntimeShield />
        <SystemResources />
        <CapabilityMatrix />
      </div>
    </RoutePage>
  );
}
