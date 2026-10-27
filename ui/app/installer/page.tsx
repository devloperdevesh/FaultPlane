import RoutePage from "@/components/layout/RoutePage";
import ManifestInstaller from "@/components/installer/ManifestInstaller";

export default function InstallerPage() {
  return (
    <RoutePage
      eyebrow="Deploy"
      title="Kubernetes Sidecar Installer"
      description="Configure FaultPlane and generate a Kubernetes manifest or Helm values file to review, copy or download."
    >
      <ManifestInstaller />
    </RoutePage>
  );
}
