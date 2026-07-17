import Panel from "@/components/ui/Panel";

export default function SyscallTrace() {
  return (
    <Panel>
      <h2 className="text-white mb-4">eBPF Syscall Trace</h2>

      <div
        className="
bg-black
rounded-lg
p-4
font-mono
text-xs
text-green-400
"
      >
        <p>[OK] kprobe tcp_connect attached</p>

        <p>[OK] socket state tracking enabled</p>

        <p>[OK] ingress packets captured</p>

        <p>[OK] overhead 0.8%</p>
      </div>
    </Panel>
  );
}
