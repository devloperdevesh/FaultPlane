"use client";

export default function GitTimeline() {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-950 p-6">
      <h2 className="mb-6 text-sm font-semibold text-white">
        Execution Timeline
      </h2>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        No execution timeline events available.
      </div>
    </div>
  );
}
