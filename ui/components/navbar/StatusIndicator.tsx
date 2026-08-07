"use client";

import { CircleCheck, CircleAlert } from "lucide-react";

interface StatusIndicatorProps {
  status?: "healthy" | "degraded";
}

export default function StatusIndicator({
  status = "healthy",
}: StatusIndicatorProps) {
  const healthy = status === "healthy";

  return (
    <div className="flex items-center gap-2 rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2">
      {healthy ? (
        <CircleCheck className="h-4 w-4 text-emerald-400" />
      ) : (
        <CircleAlert className="h-4 w-4 text-amber-400" />
      )}

      <div className="flex flex-col">
        <span className="text-xs text-zinc-500">Gateway</span>

        <span
          className={
            healthy
              ? "text-sm font-medium text-emerald-400"
              : "text-sm font-medium text-amber-400"
          }
        >
          {healthy ? "Healthy" : "Degraded"}
        </span>
      </div>
    </div>
  );
}
