const metrics = [
  ["Execution", "42ms"],

  ["Memory", "84MB"],

  ["Failures", "0"],
];

export default function SandboxMetrics() {
  return (
    <div className="grid gap-4 md:grid-cols-3">
      {metrics.map(([k, v]) => (
        <div
          key={k}
          className="
    rounded-xl
    border
    border-white/10
    bg-zinc-900
    p-4
    "
        >
          <p className="text-xs text-zinc-500">{k}</p>

          <p className="mt-2 text-white">{v}</p>
        </div>
      ))}
    </div>
  );
}
