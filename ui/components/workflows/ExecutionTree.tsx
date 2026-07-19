"use client";

import BranchNode from "./BranchNode";

const branches = [
  {
    name: "Primary Execution",
    status: "running",
    confidence: 96,
  },

  {
    name: "Speculative Branch A",
    status: "speculative",
    confidence: 82,
  },

  {
    name: "Speculative Branch B",
    status: "speculative",
    confidence: 74,
  },

  {
    name: "Rollback Branch",
    status: "failed",
    confidence: 31,
  },
] as const;

export default function ExecutionTree() {
  return (
    <div
      className="
space-y-6
"
    >
      <div
        className="
flex
items-center
justify-center
gap-6
flex-wrap
"
      >
        {branches.map((branch) => (
          <BranchNode
            key={branch.name}
            name={branch.name}
            status={branch.status}
            confidence={branch.confidence}
          />
        ))}
      </div>
    </div>
  );
}
