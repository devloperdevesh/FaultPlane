import RuntimeHealth from "./RuntimeHealth";
import KernelStatus from "./KernelStatus";

export default function SystemResources() {
  return (
    <div className="space-y-5">
      <RuntimeHealth />

      <KernelStatus />
    </div>
  );
}
