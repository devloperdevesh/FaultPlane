"use client";

import { motion } from "framer-motion";

export default function HijackDetection() {
  return (
    <div
      className="
rounded-xl
border
border-red-500/20
bg-red-500/5
p-6
"
    >
      <h3
        className="
text-sm
font-semibold
text-white
"
      >
        Hijack Detection
      </h3>

      <motion.div
        animate={{
          opacity: [0.4, 1, 0.4],
        }}
        transition={{
          repeat: Infinity,
          duration: 2,
        }}
        className="
mt-5
text-3xl
font-mono
text-red-400
"
      >
        0
      </motion.div>

      <p
        className="
mt-2
text-xs
text-zinc-400
"
      >
        No unauthorized transport interception detected
      </p>
    </div>
  );
}
