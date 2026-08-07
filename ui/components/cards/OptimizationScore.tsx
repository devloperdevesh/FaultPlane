"use client";

import { motion } from "framer-motion";

export default function OptimizationScore() {
  return (
    <div
      className="
rounded-xl
border
border-white/10
bg-zinc-950
p-6
"
    >
      <h3
        className="
text-sm
text-white
font-semibold
"
      >
        Optimization Score
      </h3>

      <div
        className="
mt-6
flex
items-center
justify-center
"
      >
        <motion.div
          animate={{
            rotate: 360,
          }}
          transition={{
            duration: 8,
            repeat: Infinity,
            ease: "linear",
          }}
          className="
h-32
w-32
rounded-full
border-8
border-emerald-500/40
flex
items-center
justify-center
"
        >
          <span
            className="
text-3xl
font-mono
text-emerald-400
"
          >
            92
          </span>
        </motion.div>
      </div>
    </div>
  );
}
