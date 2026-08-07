"use client";

import clsx from "clsx";

interface Props {
  status: "ACTIVE" | "RECOVERING" | "FAILED";
}

export default function WorkerStatus({ status }: Props) {
  return (
    <span
      className={clsx(
        "rounded-full px-3 py-1 text-xs font-medium",
        status === "ACTIVE" && "bg-emerald-500/10 text-emerald-400",

        status === "RECOVERING" && "bg-yellow-500/10 text-yellow-400",

        status === "FAILED" && "bg-red-500/10 text-red-400",
      )}
    >
      {status}
    </span>
  );
}
