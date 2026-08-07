export default function TableSkeleton() {
  return (
    <div className="space-y-3 animate-pulse">
      {Array.from({ length: 5 }).map((_, i) => (
        <div
          key={i}
          className="
    h-10
    rounded-lg
    bg-zinc-900
    "
        />
      ))}
    </div>
  );
}
