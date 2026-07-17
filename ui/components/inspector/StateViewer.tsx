const states = [
  "Memory snapshot loaded",
  "Tool execution state active",
  "Context window synchronized",
];

export default function StateViewer() {
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
      <h2 className="text-sm font-semibold text-white">Runtime State</h2>

      <div className="mt-4 space-y-3">
        {states.map((state) => (
          <div
            key={state}
            className="
              rounded-lg
              bg-zinc-900
              px-3
              py-2
              text-xs
              text-zinc-300
              "
          >
            {state}
          </div>
        ))}
      </div>
    </div>
  );
}
