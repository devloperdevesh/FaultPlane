export default function CostCard() {
  return (
    <div
      className="
        rounded-xl
        border
        border-zinc-800
        bg-zinc-950
        p-5
        "
    >
      <p className="text-xs uppercase text-zinc-500">Cost Insights</p>

      <h3 className="mt-2 text-2xl font-semibold text-white">$42.50</h3>

      <p className="mt-2 text-xs text-zinc-500">
        Estimated compute optimization
      </p>
    </div>
  );
}
