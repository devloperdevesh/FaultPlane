"use client";

import { motion } from "framer-motion";

const stages = [
  "Failure Detected",
  "Impact Propagating",
  "Containment Active",
  "Recovery Complete",
];

export default function FailurePropagation() {
  return (
    <div
      className="
rounded-xl
border
border-zinc-800
bg-zinc-950
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
        Failure Propagation
      </h3>

      <div
        className="
mt-6
space-y-4
"
      >
        {stages.map((stage, index) => (
          <motion.div
            key={stage}
            initial={{
              opacity: 0,
              x: -20,
            }}
            animate={{
              opacity: 1,
              x: 0,
            }}
            transition={{
              delay: index * 0.2,
            }}
            className="
flex
items-center
gap-4
"
          >
            <motion.div
              animate={{
                scale: [1, 1.3, 1],
              }}
              transition={{
                repeat: Infinity,
                duration: 2,
              }}
              className="
h-3
w-3
rounded-full
bg-red-500
"
            />

            <span
              className="
text-sm
text-zinc-300
"
            >
              {stage}
            </span>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
