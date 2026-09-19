"use client";

export default function CostSavings() {
  return (
    <section className="rounded-xl border border-white/10 bg-zinc-950 p-6">
      <div className="mb-5">
        <h2 className="text-sm font-semibold text-white">
          Cost Savings
        </h2>
        <p className="mt-1 text-xs text-zinc-500">
          Runtime cost and resource optimization visibility
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        No cost savings data available.
      </div>
    </section>
  );
}
