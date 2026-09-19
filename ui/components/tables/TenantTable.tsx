"use client";

export default function TenantTable() {
  return (
    <section className="rounded-2xl border border-white/10 bg-zinc-950/80 p-6 backdrop-blur-xl">
      <div className="mb-6">
        <h2 className="text-lg font-semibold text-white">
          Tenant Resource Overview
        </h2>

        <p className="mt-1 text-sm text-zinc-500">
          Multi-tenant isolation and quota governance
        </p>
      </div>

      <div className="rounded-xl border border-white/10 bg-zinc-900/60 p-4 text-sm text-zinc-500">
        Tenant resource data is not available from the current backend.
      </div>
    </section>
  );
}
