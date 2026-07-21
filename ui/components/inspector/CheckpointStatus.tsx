export default function CheckpointStatus() {
  return (
    <div
      className="
    rounded-xl
    border
    border-white/10
    bg-zinc-900
    p-4
    "
    >
      <p className="text-xs text-zinc-500">Last checkpoint</p>

      <p className="mt-2 text-white font-mono">12s ago</p>
    </div>
  );
}
