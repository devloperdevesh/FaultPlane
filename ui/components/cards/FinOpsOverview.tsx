export default function FinOpsOverview() {
  return (
    <section className="rounded-xl border border-white/10 bg-zinc-950 p-6">
      <div className="mb-5">
        <h2 className="text-sm font-semibold text-white">FinOps</h2>
        <p className="mt-1 text-xs text-zinc-500">
          Experimental cost and resource optimization visibility.
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-5 text-sm text-zinc-500">
        Live FinOps, billing, revenue, savings, and customer cost data are not
        exposed by the current backend.
      </div>
    </section>
  );
}
