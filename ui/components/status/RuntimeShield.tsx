"use client";

export default function RuntimeShield() {
  return (
    <div className="rounded-2xl border border-white/10 bg-zinc-950/80 p-6 backdrop-blur-xl">
      <div>
        <h2 className="text-lg font-semibold text-white">
          Runtime Shield
        </h2>

        <p className="mt-1 text-sm text-zinc-500">
          Runtime security and workload protection visibility
        </p>
      </div>

      <div className="mt-6 rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        Runtime security metrics are not reported by the current backend.
      </div>
    </div>
  );
}
