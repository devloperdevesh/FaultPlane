"use client";

import CommitNode from "./CommitNode";
import BranchGraph from "./BranchGraph";

const checkpoints = [
  {
    id: "hcp_v39",
    type: "checkpoint",
    time: "10:30:12",
  },

  {
    id: "hcp_v40",
    type: "checkpoint",
    time: "10:35:20",
  },

  {
    id: "hcp_v41",
    type: "checkpoint",
    time: "10:40:02",
  },

  {
    id: "rollback",
    type: "rollback",
    time: "10:42:18",
  },

  {
    id: "hcp_v42",
    type: "checkpoint",
    time: "10:43:01",
  },
] as const;

export default function CheckpointTimeline() {
  return (
    <section
      className="
rounded-xl
border
border-zinc-800
bg-zinc-950
p-6
"
    >
      <div className="mb-8">
        <h2
          className="
text-sm
font-semibold
text-white
"
        >
          Checkpoint History Timeline
        </h2>

        <p
          className="
mt-1
text-xs
text-zinc-500
"
        >
          Agent state snapshots and recovery lineage
        </p>
      </div>

      <div
        className="
flex
items-start
overflow-x-auto
"
      >
        {checkpoints.map((item, index) => (
          <div
            key={item.id}
            className="
flex
items-center
"
          >
            <CommitNode id={item.id} type={item.type} time={item.time} />

            {index !== checkpoints.length - 1 && (
              <div
                className="
mx-6
h-px
w-20
bg-zinc-700
"
              />
            )}
          </div>
        ))}
      </div>

      <BranchGraph />
    </section>
  );
}
