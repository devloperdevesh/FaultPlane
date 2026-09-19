"use client";

export default function StateViewer() {
  return (
    <section className="rounded-xl border border-white/10 bg-zinc-950/80 p-6">
      <div className="mb-5">
        <h2 className="text-sm font-semibold text-white">
          Runtime State
        </h2>
        <p className="mt-1 text-xs text-zinc-500">
          Current runtime state and recovery context
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        No runtime state data available.
      </div>
    </section>
  );
}
