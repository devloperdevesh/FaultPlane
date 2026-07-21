import HookStatus from "./HookStatus";
import KernelEvents from "./KernelEvents";
import SocketTelemetry from "./SocketTelemetry";
import SyscallTrace from "./SyscallTrace";

export default function EbpfMonitor() {
  return (
    <div className="space-y-6">
      <HookStatus />

      <div className="grid gap-6 xl:grid-cols-2">
        <SocketTelemetry />

        <KernelEvents />
      </div>

      <SyscallTrace />
    </div>
  );
}
