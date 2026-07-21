"use client";

import { motion } from "framer-motion";

interface WorkerImpactProps {
  name: string;

  status: "healthy" | "degraded" | "failed" | "recovering";
}

const statusConfig = {
  healthy: {
    label: "Healthy",
    color: "border-emerald-500/40 bg-emerald-500/10 text-emerald-400",
    dot: "bg-emerald-400",
  },

  degraded: {
    label: "Degraded",
    color: "border-orange-500/40 bg-orange-500/10 text-orange-400",
    dot: "bg-orange-400",
  },

  failed: {
    label: "Failed",
    color: "border-red-500/40 bg-red-500/10 text-red-400",
    dot: "bg-red-400",
  },

  recovering: {
    label: "Recovering",
    color: "border-blue-500/40 bg-blue-500/10 text-blue-400",
    dot: "bg-blue-400",
  },
};

export default function WorkerImpact({ name, status }: WorkerImpactProps) {
  const config = statusConfig[status];

  return (
    <motion.div
      whileHover={{
        y: -4,
        scale: 1.02,
      }}
      transition={{
        duration: 0.2,
      }}
      className="
      rounded-xl
      border
      border-white/10
      bg-zinc-950/80
      p-5
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
        <h3
          className="
          text-sm
          font-medium
          text-white
          "
        >
          {name}
        </h3>

        <span
          className={`
          h-2.5
          w-2.5
          rounded-full
          ${config.dot}
          `}
        />
      </div>

      <div
        className={`
        mt-5
        inline-flex
        rounded-full
        border
        px-3
        py-1
        text-xs
        font-medium
        ${config.color}
        `}
      >
        {config.label}
      </div>

      <div
        className="
        mt-5
        space-y-2
        text-xs
        text-zinc-500
        "
      >
        <div
          className="
          flex
          justify-between
          "
        >
          <span>Impact</span>

          <span className="text-zinc-300">
            {status === "failed"
              ? "Critical"
              : status === "degraded"
                ? "Medium"
                : status === "recovering"
                  ? "Restoring"
                  : "None"}
          </span>
        </div>

        <div
          className="
          flex
          justify-between
          "
        >
          <span>Recovery</span>

          <span className="text-zinc-300">
            {status === "recovering"
              ? "In Progress"
              : status === "failed"
                ? "Pending"
                : "Stable"}
          </span>
        </div>
      </div>
    </motion.div>
  );
}
