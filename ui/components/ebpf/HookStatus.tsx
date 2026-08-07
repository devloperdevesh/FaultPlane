import Panel from "@/components/ui/Panel";

const hooks = [
  {
    name: "XDP Ingress Hook",
    status: "ACTIVE",
  },
  {
    name: "TC Socket Monitor",
    status: "ACTIVE",
  },
  {
    name: "Kprobe Tracepoint",
    status: "ACTIVE",
  },
];

export default function HookStatus() {
  return (
    <Panel>
      <h2 className="text-white font-semibold mb-4">eBPF Kernel Hooks</h2>

      <div className="space-y-3">
        {hooks.map((hook) => (
          <div
            key={hook.name}
            className="
flex
justify-between
rounded-lg
bg-zinc-900
p-3
"
          >
            <span className="text-zinc-300">{hook.name}</span>

            <span
              className="
text-green-400
text-sm
font-mono
"
            >
              ● {hook.status}
            </span>
          </div>
        ))}
      </div>
    </Panel>
  );
}
