"use client";

export default function InterceptorLogs() {
  return (
    <div className="rounded-2xl border border-white/10 bg-zinc-950/80 p-6 backdrop-blur-xl">
      <h2 className="text-lg font-semibold text-white">
        Transport Interceptor
      </h2>

      <p className="mt-1 text-sm text-zinc-500">
        Network interception and routing decisions
      </p>

      <div className="mt-6 rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        No interceptor events available.
      </div>
    </div>
  );
}
