"use client";

import SnapshotTimeline from "./SnapshotTimeline";
import DiffBlock from "./DiffBlock";

type ChangeType = "modified" | "removed" | "added";

interface StateChange {
  type: ChangeType;
  label: string;
  before: string;
  after: string;
}

const changes: StateChange[] = [
  {
    type: "modified",
    label: "memory",
    before: "245MB",
    after: "260MB",
  },

  {
    type: "removed",
    label: "user_context",
    before: "old value",
    after: "",
  },

  {
    type: "added",
    label: "tools",
    before: "",
    after: "search_enabled",
  },
];

export default function StateDiffViewer() {
  return (
    <section
      className="
      space-y-6
      rounded-2xl
      border
      border-white/10
      bg-zinc-950/80
      p-8
      "
    >
      <div>
        <h2
          className="
          text-lg
          font-semibold
          text-white
          "
        >
          Time Travel Debugger
        </h2>

        <p
          className="
          mt-1
          text-sm
          text-zinc-500
          "
        >
          Compare agent state between checkpoints
        </p>
      </div>

      <SnapshotTimeline />

      <div
        className="
        grid
        gap-4
        "
      >
        {changes.map((change) => (
          <DiffBlock
            key={change.label}
            type={change.type}
            label={change.label}
            before={change.before}
            after={change.after}
          />
        ))}
      </div>
    </section>
  );
}
