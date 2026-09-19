"use client";

import { useEffect, useState } from "react";
import { CircleAlert, CircleCheck, CircleDashed } from "lucide-react";

import { getHealth } from "@/lib/api";

type GatewayStatus = "healthy" | "degraded" | "unavailable";

export default function StatusIndicator() {
  const [status, setStatus] = useState<GatewayStatus>("unavailable");

  useEffect(() => {
    let mounted = true;

    const checkHealth = async () => {
      try {
        const health = await getHealth();

        if (!mounted) return;

        setStatus(health.status === "healthy" ? "healthy" : "degraded");
      } catch {
        if (mounted) {
          setStatus("unavailable");
        }
      }
    };

    void checkHealth();

    const interval = window.setInterval(checkHealth, 5000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  const config = {
    healthy: {
      icon: CircleCheck,
      iconClass: "text-emerald-400",
      textClass: "text-emerald-400",
      label: "Healthy",
    },
    degraded: {
      icon: CircleAlert,
      iconClass: "text-amber-400",
      textClass: "text-amber-400",
      label: "Degraded",
    },
    unavailable: {
      icon: CircleDashed,
      iconClass: "text-zinc-500",
      textClass: "text-zinc-400",
      label: "Unavailable",
    },
  }[status];

  const Icon = config.icon;

  return (
    <div className="flex items-center gap-2 rounded-lg border border-zinc-800 bg-zinc-900 px-3 py-2">
      <Icon className={`h-4 w-4 ${config.iconClass}`} />

      <div className="flex flex-col">
        <span className="text-xs text-zinc-500">Gateway</span>

        <span className={`text-sm font-medium ${config.textClass}`}>
          {config.label}
        </span>
      </div>
    </div>
  );
}
