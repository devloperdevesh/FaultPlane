"use client";

import { ReactNode } from "react";
import clsx from "clsx";

interface SidebarSectionProps {
  title: string;
  children: ReactNode;
  className?: string;
}

export default function SidebarSection({
  title,
  children,
  className,
}: SidebarSectionProps) {
  return (
    <section
      className={clsx(
        "mb-8 animate-in fade-in slide-in-from-left-2 duration-300",
        className,
      )}
    >
      <h2
        className="
          mb-3
          px-4
          text-[11px]
          font-semibold
          uppercase
          tracking-[0.22em]
          text-zinc-500
        "
      >
        {title}
      </h2>

      <div className="space-y-1">{children}</div>
    </section>
  );
}
