"use client";

import { motion } from "framer-motion";

import VariableDiffRow from "./VariableDiff";
import StateSummary from "./StateSummary";

import type { VariableDiff } from "@/lib/types";

const changes: VariableDiff[] = [
  {
    key: "model",
    before: "gpt-4o",
    after: "claude-4",
    type: "modified",
  },

  {
    key: "memory_tokens",
    before: "4200",
    after: "5100",
    type: "modified",
  },

  {
    key: "checkpoint",
    before: "cp-1024",
    after: "restored",
    type: "added",
  },

  {
    key: "tools",
    before: "enabled",
    after: "enabled",
    type: "unchanged",
  },
];

export default function StateDiff() {
  return (
    <section
      className="
        rounded-2xl
        border
        border-white/10
        bg-zinc-950/80
        p-8
        space-y-8
      "
    >
      {/* Header */}

      <div>
        <h2
          className="
            text-lg
            font-semibold
            text-white
          "
        >
          Checkpoint State Diff
        </h2>

        <p
          className="
            mt-1
            text-sm
            text-zinc-500
          "
        >
          Compare runtime state against recovery snapshot
        </p>
      </div>

      {/* Summary */}

      <StateSummary />

      {/* Diff List */}

      <div className="space-y-4">
        {changes.map((change, index) => (
          <motion.div
            key={change.key}
            initial={{
              opacity: 0,
              y: 20,
            }}
            animate={{
              opacity: 1,
              y: 0,
            }}
            transition={{
              delay: index * 0.08,
              duration: 0.25,
            }}
          >
            <VariableDiffRow data={change} />
          </motion.div>
        ))}
      </div>
    </section>
  );
}
