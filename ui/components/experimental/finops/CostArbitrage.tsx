import CostComparison from "./CostComparison";
import SavingsCard from "./SavingsCard";

export default function CostArbitrage() {
  return (
    <section
      className="
space-y-5
rounded-2xl
border
border-white/10
bg-zinc-950
p-6
"
    >
      <h2 className="text-white font-semibold">Cloud Runtime Cost</h2>

      <div className="grid gap-5 xl:grid-cols-3">
        <div className="rounded-xl bg-zinc-900 p-5">
          <p className="text-xs text-zinc-500">Standard Compute</p>

          <p className="mt-2 text-white">$4.20/hr</p>
        </div>

        <div className="rounded-xl bg-zinc-900 p-5">
          <p className="text-xs text-zinc-500">Optimized Runtime</p>

          <p className="mt-2 text-emerald-400">$0.85/hr</p>
        </div>

        <SavingsCard />
      </div>

      <CostComparison />
    </section>
  );
}
