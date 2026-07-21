export default function RecoveryCard() {
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
      <p className="text-xs uppercase text-zinc-500">Recovery Events</p>

      <h3 className="mt-2 text-2xl font-semibold text-white">24</h3>

      <p className="mt-2 text-xs text-zinc-500">
        Successful failover recoveries
      </p>
    </div>
  );
}
