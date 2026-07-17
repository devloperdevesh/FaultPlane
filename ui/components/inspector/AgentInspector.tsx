export default function AgentInspector() {
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
      <h2 className="text-sm font-semibold text-white">Agent Inspector</h2>

      <p className="mt-1 text-xs text-zinc-500">Runtime execution details</p>

      <div className="mt-5 space-y-3">
        <div className="flex justify-between">
          <span className="text-sm text-zinc-400">Agent ID</span>

          <span className="text-sm text-white">agent-worker-01</span>
        </div>

        <div className="flex justify-between">
          <span className="text-sm text-zinc-400">Status</span>

          <span className="text-sm text-emerald-400">Running</span>
        </div>

        <div className="flex justify-between">
          <span className="text-sm text-zinc-400">Context</span>

          <span className="text-sm text-white">84k tokens</span>
        </div>
      </div>
    </div>
  );
}
