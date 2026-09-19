"use client";

export default function CostArbitrage() {
  return (
    <section className="space-y-5 rounded-2xl border border-white/10 bg-zinc-950 p-6">
      <div>
        <h2 className="font-semibold text-white">Cloud Runtime Cost</h2>
        <p className="mt-1 text-xs text-zinc-500">
          Experimental FinOps view. No live billing or cloud-cost telemetry is
          connected to the current runtime API.
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-5 text-sm text-zinc-500">
        Standard compute pricing, optimized runtime pricing, savings estimates,
        and billing data are not currently reported by FaultPlane.
      </div>
    </section>
  );
}
