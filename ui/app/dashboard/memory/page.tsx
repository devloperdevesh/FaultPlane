"use client";

import RoutePage from "@/components/layout/RoutePage";
import MemoryCockpitOverview from "@/components/memory/MemoryCockpitOverview";

export default function MemoryCockpitPage() {
  return (
    <RoutePage
      eyebrow="Buffer Architecture & Isolation Telemetry"
      title="Cache-Line Memory Matrix"
      description="64-byte aligned multi-tenant buffer pool with CAS slot synchronization, runtime GC pause telemetry, and logical domain fencing (software isolation mode)."
    >
      <MemoryCockpitOverview />
    </RoutePage>
  );
}
