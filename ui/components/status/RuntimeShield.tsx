"use client";

import { motion } from "framer-motion";

export default function RuntimeShield() {
  return (
    <div
      className="
rounded-2xl
border
border-emerald-500/20
bg-zinc-950/80
p-6
backdrop-blur-xl
"
    >
      <div
        className="
flex
items-center
justify-between
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
            Runtime Shield
          </h2>

          <p
            className="
mt-1
text-sm
text-zinc-500
"
          >
            AI workload security posture
          </p>
        </div>

        <motion.div
          animate={{
            scale: [1, 1.15, 1],
          }}
          transition={{
            duration: 2,
            repeat: Infinity,
          }}
          className="
h-4
w-4
rounded-full
bg-emerald-400
"
        />
      </div>

      <div
        className="
mt-6
grid
grid-cols-3
gap-4
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
            Sandbox
          </p>

          <p
            className="
mt-2
text-white
font-mono
"
          >
            48 Active
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
            Blocked
          </p>

          <p
            className="
mt-2
text-red-400
font-mono
"
          >
            12 Events
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
            Risk
          </p>

          <p
            className="
mt-2
text-emerald-400
font-mono
"
          >
            LOW
          </p>
        </div>
      </div>
    </div>
  );
}
