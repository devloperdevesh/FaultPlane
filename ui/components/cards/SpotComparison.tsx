"use client";

import { motion } from "framer-motion";

export default function SpotComparison() {
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
font-semibold
text-white
"
      >
        Spot Arbitrage
      </h3>

      <div
        className="
mt-6
grid
grid-cols-3
items-center
gap-4
"
      >
        <div
          className="
rounded-xl
bg-red-500/10
p-4
text-center
"
        >
          <p
            className="
text-xs
text-zinc-400
"
          >
            On Demand
          </p>

          <p
            className="
mt-2
text-xl
font-mono
text-red-400
"
          >
            $0.42
          </p>
        </div>

        <motion.div
          animate={{
            x: [0, 20, 0],
          }}
          transition={{
            repeat: Infinity,
            duration: 2,
          }}
          className="
text-xl
text-emerald-400
"
        >
          →
        </motion.div>

        <div
          className="
rounded-xl
bg-emerald-500/10
p-4
text-center
"
        >
          <p
            className="
text-xs
text-zinc-400
"
          >
            Spot
          </p>

          <p
            className="
mt-2
text-xl
font-mono
text-emerald-400
"
          >
            $0.13
          </p>
        </div>
      </div>
    </div>
  );
}
