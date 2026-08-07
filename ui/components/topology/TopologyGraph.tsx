import GatewayNode from "./GatewayNode";
import WorkerNode from "./WorkerNode";
import Connection from "./Connection";

import Node from "./Node";

export default function TopologyGraph() {
  return (
    <div
      className="
      rounded-xl
      border
      border-zinc-800
      bg-zinc-950
      p-8
      "
    >
      <div className="mb-6">
        <h2 className="text-sm font-semibold text-white">
          Infrastructure Topology
        </h2>

        <p className="text-xs text-zinc-500">Runtime dependency graph</p>
      </div>

      <div
        className="
        flex
        flex-col
        items-center
        "
      >
        <GatewayNode name="Gateway" status="Healthy" />

        <Connection label="TCP Stream" />

        <WorkerNode id="Worker-01" status="Running" />
        <Node label="Client" type="External" />

        <Connection label="Checkpoint Sync" />

        <div
          className="
          rounded-xl
          border
          border-emerald-500/30
          bg-zinc-950
          px-6
          py-4
          text-center
          "
        >
          <p className="text-sm font-semibold text-white">Checkpoint Store</p>

          <p className="mt-2 text-xs text-emerald-400">
            State Persistence Active
          </p>
        </div>
      </div>
    </div>
  );
}
