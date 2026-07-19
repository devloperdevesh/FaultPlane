"use client";

import { motion } from "framer-motion";

interface Props {
  name: string;
  status: "healthy" | "warning" | "critical" | "recovered";
}

const colors = {
  healthy: "border-emerald-500 text-emerald-400",
  warning: "border-yellow-500 text-yellow-400",
  critical: "border-red-500 text-red-400",
  recovered: "border-blue-500 text-blue-400",
};

export default function ServiceNode({ name, status }: Props) {
  return (
    <motion.div
      initial={{
        scale: 0.8,
        opacity: 0,
      }}
      animate={{
        scale: 1,
        opacity: 1,
      }}
      transition={{
        duration: 0.4,
      }}
      className={`
rounded-xl
border
bg-zinc-950
px-6
py-4
text-center
${colors[status]}
`}
    >
      <p className="text-sm">{name}</p>

      <p className="text-xs mt-2 opacity-70">{status}</p>
    </motion.div>
  );
}
