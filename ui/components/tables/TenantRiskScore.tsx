export default function TenantRiskScore() {
  return (
    <div className="rounded-xl border border-white/10 bg-zinc-950 p-5">
      <p className="text-xs text-zinc-500">
        Isolation Risk Score
      </p>

      <p className="mt-3 text-sm text-zinc-300">
        Unavailable
      </p>

      <p className="mt-2 text-xs leading-5 text-zinc-500">
        Live tenant isolation risk scoring is not exposed by the current
        runtime API.
      </p>
    </div>
  );
}
