"use client";

import { motion } from "framer-motion";

export default function TenantRiskScore() {
  return (
    <div
      className="
rounded-xl
border
border-white/10
bg-zinc-950
p-5
"
    >
      <p
        className="
text-xs
text-zinc-500
"
      >
        Isolation Risk Score
      </p>

      <motion.p
        initial={{
          scale: 0.8,
        }}
        animate={{
          scale: 1,
        }}
        className="
mt-3
text-3xl
font-mono
text-emerald-400
"
      >
        12%
      </motion.p>

      <p
        className="
mt-2
text-xs
text-zinc-400
"
      >
        Low resource contention detected
      </p>
    </div>
  );
}
