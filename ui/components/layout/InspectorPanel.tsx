"use client";

import type { ReactNode } from "react";
import { motion } from "framer-motion";

import AgentInspector from "@/components/inspector/AgentInspector";

interface InspectorPanelProps {
  children?: ReactNode;
}

export default function InspectorPanel({ children }: InspectorPanelProps) {
  return (
    <motion.aside
      initial={{
        opacity: 0,
        x: 40,
      }}
      animate={{
        opacity: 1,
        x: 0,
      }}
      transition={{
        duration: 0.25,
      }}
      className="
      flex
      h-full
      w-80
      shrink-0
      flex-col
      border-l
      border-zinc-800
      bg-zinc-950/80
      backdrop-blur-xl
      "
    >
      {/* Header */}

      <div
        className="
        border-b
        border-zinc-800
        px-5
        py-4
        "
      >
        <h2
          className="
          text-sm
          font-semibold
          text-white
          "
        >
          Runtime Inspector
        </h2>

        <p
          className="
          mt-1
          text-xs
          text-zinc-500
          "
        >
          Live agent context and state
        </p>
      </div>

      {/* Content */}

      <div
        className="
        flex-1
        overflow-y-auto
        p-5
        "
      >
        {children ?? <AgentInspector />}
      </div>
    </motion.aside>
  );
}
