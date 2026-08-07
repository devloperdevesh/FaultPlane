export default function CheckpointDetails() {
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
      <h2 className="text-sm font-semibold text-white">Checkpoint Details</h2>

      <div className="mt-4 space-y-3">
        <div className="flex justify-between">
          <span className="text-zinc-400 text-sm">ID</span>

          <span className="text-white text-sm">cp-1024</span>
        </div>

        <div className="flex justify-between">
          <span className="text-zinc-400 text-sm">Size</span>

          <span className="text-white text-sm">18 MB</span>
        </div>

        <div className="flex justify-between">
          <span className="text-zinc-400 text-sm">Recovery</span>

          <span className="text-emerald-400 text-sm">Ready</span>
        </div>
      </div>
    </div>
  );
}
