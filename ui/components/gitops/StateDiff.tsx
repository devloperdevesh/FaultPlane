"use client";

export default function StateDiff() {
  return (
    <section className="space-y-8 rounded-2xl border border-white/10 bg-zinc-950/80 p-8">
      <div>
        <h2 className="text-lg font-semibold text-white">
          Checkpoint State Diff
        </h2>

        <p className="mt-1 text-sm text-zinc-500">
          Compare runtime state against recovery snapshot
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        Checkpoint state diff data is not available.
      </div>
    </section>
  );
}
