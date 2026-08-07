export default function CardSkeleton() {
  return (
    <div
      className="
    animate-pulse
    rounded-2xl
    border
    border-white/10
    bg-zinc-900
    p-6
    "
    >
      <div className="h-4 w-24 rounded bg-zinc-700" />

      <div className="mt-4 h-8 w-32 rounded bg-zinc-800" />
    </div>
  );
}
