"use client";

import RoutePage from "@/components/layout/RoutePage";
import MemoryCockpitOverview from "@/components/memory/MemoryCockpitOverview";

export default function MemoryCockpitPage() {
  return (
    <RoutePage
      eyebrow="Silicon Architecture & Memory Protection"
      title="Sub-Kernel Memory Matrix"
      description="Hardware-isolated multi-tenant buffer pools, 64-byte structural cache-line alignment guard, and real-time 3D WebGL silicon fencing cockpit."
    >
      <MemoryCockpitOverview />
    </RoutePage>
  );
}
