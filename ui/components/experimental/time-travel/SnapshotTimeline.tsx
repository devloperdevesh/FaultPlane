"use client";

export default function SnapshotTimeline() {
  return (
    <div
      className="
      rounded-2xl
      border
      border-white/10
      bg-zinc-950
      p-6
      "
    >
      <h3 className="text-sm font-semibold text-white">
        Checkpoint History
      </h3>

      <div className="mt-6 rounded-xl border border-white/10 bg-zinc-900/60 p-5">
        <p className="text-sm text-zinc-300">
          Checkpoint history unavailable
        </p>

        <p className="mt-2 text-xs leading-5 text-zinc-500">
          Historical checkpoint listings are not currently exposed by the
          runtime API.
        </p>
      </div>
    </div>
  );
}
