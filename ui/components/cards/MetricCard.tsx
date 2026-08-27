"use client";

import { motion } from "framer-motion";
import type { LucideIcon } from "lucide-react";

interface Props {
  title: string;
  value: string;
  description?: string;
  icon: LucideIcon;
  status?: "healthy" | "warning" | "critical";
  trend?: string;
}

export default function MetricCard({
  title,
  value,
  description,
  icon: Icon,
  status = "healthy",
  trend,
}: Props) {
  const statusStyles = {
    healthy: "text-emerald-400",
    warning: "text-yellow-400",
    critical: "text-red-400",
  };

  const dotStyles = {
    healthy: "bg-emerald-400",
    warning: "bg-yellow-400",
    critical: "bg-red-400",
  };

  return (
    <motion.div
      whileHover={{ y: -3 }}
      transition={{ duration: 0.2 }}
      className="relative overflow-hidden rounded-2xl border border-white/10 bg-zinc-950/70 p-5 shadow-xl backdrop-blur-xl"
    >
      <div className="flex items-center justify-between">
        <div className="rounded-xl bg-white/5 p-3">
          <Icon className={statusStyles[status]} size={22} />
        </div>

        <div className={`h-2 w-2 rounded-full ${dotStyles[status]}`} />
      </div>

      <p className="mt-5 text-xs uppercase tracking-widest text-zinc-500">
        {title}
      </p>

      <h2 className="mt-2 text-3xl font-semibold text-white">{value}</h2>

      {description && (
        <p className="mt-2 text-xs text-zinc-500">{description}</p>
      )}

      {trend && (
        <p className="mt-4 text-xs text-zinc-400">{trend}</p>
      )}
    </motion.div>
  );
}
