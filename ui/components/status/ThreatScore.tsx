"use client";

import { motion } from "framer-motion";

export default function ThreatScore() {
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
      <p
        className="
text-xs
text-zinc-500
"
      >
        Threat Score
      </p>

      <motion.div
        initial={{
          width: 0,
        }}
        animate={{
          width: "18%",
        }}
        transition={{
          duration: 1,
        }}
        className="
mt-5
h-3
rounded-full
bg-emerald-500
"
      />

      <p
        className="
mt-3
font-mono
text-2xl
text-emerald-400
"
      >
        18 / 100
      </p>
    </div>
  );
}
