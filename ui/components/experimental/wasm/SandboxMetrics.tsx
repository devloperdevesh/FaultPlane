export default function SandboxMetrics() {
  return (
    <div className="rounded-xl border border-white/10 bg-zinc-900 p-4">
      <p className="text-xs text-zinc-500">Sandbox metrics</p>
      <p className="mt-2 text-sm text-zinc-300">Unavailable</p>
      <p className="mt-1 text-xs leading-5 text-zinc-500">
        Live WASM execution metrics are not exposed by the runtime API.
      </p>
    </div>
  );
}
