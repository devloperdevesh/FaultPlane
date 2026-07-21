"use client";

import { motion } from "framer-motion";

interface Props {
  label: string;
  value: number;
}

export default function QuotaProgress({ label, value }: Props) {
  return (
    <div className="space-y-2">
      <div
        className="
flex
justify-between
text-xs
text-zinc-400
"
      >
        <span>{label}</span>

        <span>{value}%</span>
      </div>

      <div
        className="
h-2
rounded-full
bg-zinc-800
overflow-hidden
"
      >
        <motion.div
          initial={{
            width: 0,
          }}
          animate={{
            width: `${value}%`,
          }}
          transition={{
            duration: 1,
            ease: "easeOut",
          }}
          className={`
h-full
rounded-full

${value >= 90 ? "bg-red-500" : value >= 70 ? "bg-orange-500" : "bg-emerald-500"}

`}
        />
      </div>
    </div>
  );
}
