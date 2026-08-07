"use client";

import { motion } from "framer-motion";

import ExecutionTree from "./ExecutionTree";

export default function SpeculativePath() {
  return (
    <section
      className="
rounded-2xl
border
border-white/10
bg-zinc-950/80
backdrop-blur-xl
p-8
"
    >
      <div className="mb-8">
        <h2
          className="
text-lg
font-semibold
text-white
"
        >
          Speculative Execution Engine
        </h2>

        <p
          className="
mt-1
text-sm
text-zinc-500
"
        >
          Parallel workflow prediction and optimal path selection
        </p>
      </div>

      <div
        className="
relative
rounded-xl
border
border-white/10
bg-black/30
p-8
"
      >
        <motion.div
          animate={{
            opacity: [0.3, 0.8, 0.3],
          }}
          transition={{
            duration: 3,
            repeat: Infinity,
          }}
          className="
absolute
inset-0
bg-gradient-to-r
from-purple-500/10
via-transparent
to-blue-500/10
"
        />

        <ExecutionTree />
      </div>

      <div
        className="
mt-6
grid
gap-4
md:grid-cols-3
"
      >
        <div
          className="
rounded-xl
bg-zinc-900
p-4
"
        >
          <p
            className="
text-xs
text-zinc-500
"
          >
            Active Branches
          </p>

          <p
            className="
mt-2
font-mono
text-white
"
          >
            3
          </p>
        </div>

        <div
          className="
rounded-xl
bg-zinc-900
p-4
"
        >
          <p
            className="
text-xs
text-zinc-500
"
          >
            Best Path
          </p>

          <p
            className="
mt-2
font-mono
text-emerald-400
"
          >
            96%
          </p>
        </div>

        <div
          className="
rounded-xl
bg-zinc-900
p-4
"
        >
          <p
            className="
text-xs
text-zinc-500
"
          >
            Rollback
          </p>

          <p
            className="
mt-2
font-mono
text-white
"
          >
            Ready
          </p>
        </div>
      </div>
    </section>
  );
}
