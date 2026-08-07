export default function SLAProtectionCard() {
  return (
    <div
      className="
        rounded-xl
        border
        border-green-500/20
        bg-zinc-950
        p-5
        "
    >
      <p className="text-xs uppercase text-zinc-500">SLA Protected</p>

      <h3 className="mt-2 text-2xl font-semibold text-green-400">$142,520</h3>

      <p className="mt-2 text-xs text-zinc-500">
        Estimated avoided downtime impact
      </p>
    </div>
  );
}
