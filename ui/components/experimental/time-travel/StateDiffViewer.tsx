"use client";

export default function StateDiffViewer() {
  return (
    <section
      className="
      space-y-6
      rounded-2xl
      border
      border-white/10
      bg-zinc-950/80
      p-8
      "
    >
      <div>
        <h2 className="text-lg font-semibold text-white">
          Time Travel Debugger
        </h2>

        <p className="mt-1 text-sm text-zinc-500">
          Compare agent state between checkpoints
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-5">
        <p className="text-sm text-zinc-300">
          State diff unavailable
        </p>

        <p className="mt-2 text-xs leading-5 text-zinc-500">
          Historical checkpoint state and diff data are not currently exposed
          by the runtime API.
        </p>
      </div>
    </section>
  );
}
