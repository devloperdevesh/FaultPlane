"use client";

import clsx from "clsx";

export type Status =
  | "healthy"
  | "running"
  | "warning"
  | "failed"
  | "experimental";

interface Props {
  status: Status;
  label: string;
}

const config = {
  healthy: {
    text: "Healthy",
    class: "border-emerald-500/30 bg-emerald-500/10 text-emerald-400",
    icon: "●",
  },

  running: {
    text: "Running",
    class: "border-blue-500/30 bg-blue-500/10 text-blue-400",
    icon: "●",
  },

  warning: {
    text: "Warning",
    class: "border-yellow-500/30 bg-yellow-500/10 text-yellow-400",
    icon: "⚠",
  },

  failed: {
    text: "Failed",
    class: "border-red-500/30 bg-red-500/10 text-red-400",
    icon: "✕",
  },

  experimental: {
    text: "Experimental",
    class: "border-purple-500/30 bg-purple-500/10 text-purple-400",
    icon: "◈",
  },
};

export default function StatusBadge({ status, label }: Props) {
  const item = config[status];

  return (
    <span
      className={clsx(
        "inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-medium",
        item.class,
      )}
    >
      <span>{item.icon}</span>

      {label || item.text}
    </span>
  );
}
