"use client";

import RuntimeState from "./RuntimeState";
import ContextGauge from "./ContextGauge";
import CheckpointStatus from "./CheckpointStatus";

export default function AgentInspector() {
  return (
    <aside
      className="
space-y-6
rounded-2xl
border
border-white/10
bg-zinc-950/80
p-6
"
    >
      <div>
        <h2
          className="
text-sm
font-semibold
text-white
"
        >
          Agent Inspector
        </h2>

        <p
          className="
mt-1
text-xs
text-zinc-500
"
        >
          Runtime execution details
        </p>
      </div>

      <div
        className="
space-y-4
"
      >
        <div>
          <p className="text-xs text-zinc-500">State</p>

          <div className="mt-2">
            <RuntimeState state="RUNNING" />
          </div>
        </div>

        <div>
          <p className="text-xs text-zinc-500">Memory</p>

          <p className="mt-2 text-xl text-white">245 MB</p>
        </div>

        <ContextGauge value={84} />

        <CheckpointStatus />
      </div>
    </aside>
  );
}
