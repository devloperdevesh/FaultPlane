"use client";

import { motion } from "framer-motion";

interface Props {
  name: string;

  status: "running" | "success" | "failed" | "speculative";

  confidence: number;
}

const styles = {
  running: "border-blue-500/40 bg-blue-500/10 text-blue-400",

  success: "border-emerald-500/40 bg-emerald-500/10 text-emerald-400",

  failed: "border-red-500/40 bg-red-500/10 text-red-400",

  speculative: "border-purple-500/40 bg-purple-500/10 text-purple-400",
};

export default function BranchNode({ name, status, confidence }: Props) {
  return (
    <motion.div
      initial={{
        opacity: 0,
        scale: 0.8,
      }}
      animate={{
        opacity: 1,
        scale: 1,
      }}
      whileHover={{
        scale: 1.05,
      }}
      className={`
rounded-xl
border
p-4
backdrop-blur-xl
${styles[status]}
`}
    >
      <p
        className="
text-sm
font-semibold
"
      >
        {name}
      </p>

      <p
        className="
mt-2
text-xs
opacity-80
"
      >
        {status}
      </p>

      <div
        className="
mt-3
text-xs
font-mono
"
      >
        Confidence {confidence}%
      </div>
    </motion.div>
  );
}
