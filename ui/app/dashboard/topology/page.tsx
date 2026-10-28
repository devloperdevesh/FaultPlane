"use client";

import { useState } from "react";
import RoutePage from "@/components/layout/RoutePage";
import TopologyDashboard from "@/components/topology/TopologyDashboard";

export default function TopologyPage() {
  const [demo, setDemo] = useState(false);
  return (
    <RoutePage
      eyebrow="Runtime observability"
      title="Topology & blast radius"
      description="Inspect reported regions, worker health, dependency exposure and recorded recovery events."
    >
      <div className="flex flex-wrap items-center justify-between gap-3">
        <p className="text-sm text-zinc-400">
          {demo
            ? "Demo replay is isolated from live API data."
            : "Read-only observations from the FaultPlane daemon."}
        </p>
        <button
          type="button"
          onClick={() => setDemo((value) => !value)}
          className="rounded-lg border border-zinc-600 px-4 py-2 text-sm text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-400"
        >
          {demo ? "Return to live telemetry" : "Open demo replay"}
        </button>
      </div>
      <TopologyDashboard key={demo ? "demo" : "live"} demo={demo} />
    </RoutePage>
  );
}
