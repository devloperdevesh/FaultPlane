"use client";

export default function CheckpointTimeline() {
  return (
    <section className="rounded-xl border border-zinc-800 bg-zinc-950 p-6">
      <div className="mb-8">
        <h2 className="text-sm font-semibold text-white">
          Checkpoint History
        </h2>
        <p className="mt-1 text-xs text-zinc-500">
          Runtime state snapshots
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        No checkpoint history available.
      </div>
    </section>
  );
}
