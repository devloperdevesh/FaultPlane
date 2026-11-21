import RuntimeHealth from "./RuntimeHealth";
import KernelStatus from "./KernelStatus";
import HardwareVectorCard from "./HardwareVectorCard";

export default function SystemResources() {
  return (
    <div className="space-y-5">
      <HardwareVectorCard />

      <RuntimeHealth />

      <KernelStatus />
    </div>
  );
}
