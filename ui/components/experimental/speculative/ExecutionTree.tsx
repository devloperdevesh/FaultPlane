"use client";

import { motion } from "framer-motion";

interface ExecutionPathProps {
  height?: number;
  animated?: boolean;
}

export default function ExecutionPath({
  height = 48,
  animated = true,
}: ExecutionPathProps) {
  return (
    <motion.div
      initial={
        animated
          ? {
              opacity: 0,
              scaleY: 0,
            }
          : undefined
      }
      animate={
        animated
          ? {
              opacity: 1,
              scaleY: 1,
            }
          : undefined
      }
      transition={{
        duration: 0.35,
        ease: "easeOut",
      }}
      className="
        flex
        origin-top
        justify-center
        py-3
      "
      aria-hidden="true"
    >
      <div
        style={{
          height,
        }}
        className="
          w-px
          rounded-full
          bg-gradient-to-b
          from-blue-500
          via-zinc-700
          to-emerald-500
        "
      />
    </motion.div>
  );
}
