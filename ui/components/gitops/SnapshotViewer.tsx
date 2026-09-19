export default function SnapshotViewer() {
  return (
    <div className="rounded-2xl border border-white/10 bg-zinc-950 p-6">
      <h2 className="text-sm font-semibold text-white">
        Runtime Snapshot
      </h2>

      <div className="mt-5 rounded-xl bg-zinc-900 p-4">
        <p className="text-sm text-zinc-300">
          Snapshot data unavailable
        </p>

        <p className="mt-2 text-xs leading-5 text-zinc-500">
          Live runtime snapshot history and state metrics are not exposed by
          the current runtime API.
        </p>
      </div>
    </div>
  );
}
