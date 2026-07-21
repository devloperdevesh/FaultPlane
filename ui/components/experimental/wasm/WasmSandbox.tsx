import RuntimeIsolation from "./RuntimeIsolation";
import SandboxMetrics from "./SandboxMetrics";

export default function WasmSandbox() {
  return (
    <section
      className="
space-y-6
rounded-2xl
border
border-white/10
bg-zinc-950
p-6
"
    >
      <div>
        <h2 className="text-white font-semibold">WASM Sandbox</h2>

        <p className="text-xs text-zinc-500">Runtime isolation layer</p>
      </div>

      <div
        className="
grid
gap-4
md:grid-cols-4
"
      >
        <div>
          <p className="text-xs text-zinc-500">Status</p>

          <p className="text-emerald-400">ACTIVE</p>
        </div>

        <div>
          <p className="text-xs text-zinc-500">Memory Limit</p>

          <p className="text-white">128MB</p>
        </div>

        <div>
          <p className="text-xs text-zinc-500">CPU Limit</p>

          <p className="text-white">2 cores</p>
        </div>

        <div>
          <p className="text-xs text-zinc-500">Isolation</p>

          <p className="text-emerald-400">Enabled</p>
        </div>
      </div>

      <RuntimeIsolation />

      <SandboxMetrics />
    </section>
  );
}
