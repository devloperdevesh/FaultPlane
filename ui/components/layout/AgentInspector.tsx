export default function AgentInspector() {
  return (
    <div className="space-y-5">
      <h2 className="text-white font-semibold">Worker Inspector</h2>

      <div className="space-y-3 text-sm">
        <p className="text-zinc-400">
          Status
          <span className="float-right text-emerald-400">Healthy</span>
        </p>

        <p className="text-zinc-400">
          Memory
          <span className="float-right text-white">1.8GB</span>
        </p>

        <p className="text-zinc-400">
          Checkpoint
          <span className="float-right text-white">hcp_42</span>
        </p>

        <p className="text-zinc-400">
          Recovery
          <span className="float-right text-white">14ms</span>
        </p>
      </div>
    </div>
  );
}
