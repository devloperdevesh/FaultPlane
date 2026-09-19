"use client";

export default function AgentInspector() {
  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-sm font-semibold text-white">
          Runtime Inspector
        </h2>

        <p className="mt-1 text-xs text-zinc-500">
          Live worker context
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4">
        <p className="text-sm text-zinc-300">
          No worker selected
        </p>

        <p className="mt-2 text-xs leading-5 text-zinc-500">
          Select a worker from the runtime view to inspect live state and
          telemetry.
        </p>
      </div>
    </div>
  );
}
