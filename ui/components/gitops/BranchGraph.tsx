export default function BranchGraph() {
  return (
    <div
      className="
    mt-8
    space-y-3
    font-mono
    text-xs
    "
    >
      <div
        className="
    flex
    items-center
    gap-3
    "
      >
        <span className="text-cyan-400">●</span>

        <span className="text-zinc-400">primary-runtime</span>
      </div>

      <div
        className="
    ml-2
    border-l
    border-zinc-700
    pl-5
    "
      >
        <div
          className="
    flex
    gap-3
    items-center
    "
        >
          <span className="text-green-400">●</span>

          <span className="text-zinc-400">fallback-recovery</span>
        </div>
      </div>
    </div>
  );
}
