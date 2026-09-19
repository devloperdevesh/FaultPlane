"use client";

import BlastGraph from "./BlastGraph";
import ImpactTimeline from "./ImpactTimeline";
import RadiusCircle from "./RadiusCircle";

export default function BlastRadius() {
  return (
    <section className="relative overflow-hidden rounded-2xl border border-white/10 bg-zinc-950/80 p-8 backdrop-blur-xl space-y-8">
      <div className="flex items-start justify-between">
        <div>
          <h2 className="text-lg font-semibold text-white">
            Recovery Blast Radius
          </h2>
          <p className="mt-1 text-sm text-zinc-500">
            Failure propagation, impact analysis and recovery containment
          </p>
        </div>

        <div className="rounded-full border border-zinc-700 bg-zinc-900 px-4 py-2 text-xs text-zinc-400">
          No active recovery data
        </div>
      </div>

      <div className="grid gap-8 xl:grid-cols-3">
        <div className="flex items-center justify-center rounded-xl border border-white/10 bg-black/40 p-8">
          <RadiusCircle />
        </div>

        <div className="rounded-xl border border-white/10 bg-black/40 p-6">
          <BlastGraph />
        </div>

        <div className="rounded-xl border border-white/10 bg-black/40 p-6">
          <ImpactTimeline />
        </div>
      </div>

      <div className="rounded-xl border border-white/10 bg-black/40 p-6 text-sm text-zinc-500">
        No worker impact records available.
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {["Affected Nodes", "Recovery Progress", "Containment", "MTTR"].map(
          (title) => (
            <div
              key={title}
              className="rounded-xl border border-white/10 bg-zinc-900/70 p-5"
            >
              <p className="text-xs uppercase tracking-wider text-zinc-500">
                {title}
              </p>
              <p className="mt-3 font-mono text-xl text-zinc-600">—</p>
            </div>
          ),
        )}
      </div>
    </section>
  );
}
