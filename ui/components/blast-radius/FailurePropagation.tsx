"use client";

export default function FailurePropagation() {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-950 p-6">
      <h3 className="text-sm font-semibold text-white">
        Failure Propagation
      </h3>

      <div className="mt-6 rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        No active failure propagation events available.
      </div>
    </div>
  );
}
