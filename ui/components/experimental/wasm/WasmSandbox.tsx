"use client";

import RuntimeIsolation from "./RuntimeIsolation";
import SandboxMetrics from "./SandboxMetrics";

export default function WasmSandbox() {
  return (
    <section className="space-y-6 rounded-2xl border border-white/10 bg-zinc-950 p-6">
      <div>
        <h2 className="font-semibold text-white">
          WASM Sandbox
        </h2>

        <p className="text-xs text-zinc-500">
          Runtime isolation layer
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        WASM sandbox runtime status is not reported by the current backend.
      </div>

      <RuntimeIsolation />

      <SandboxMetrics />
    </section>
  );
}
