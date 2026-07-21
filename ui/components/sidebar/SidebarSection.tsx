"use client";

import type { ReactNode } from "react";
import type { LucideIcon } from "lucide-react";
import clsx from "clsx";

interface SidebarSectionProps {
  title: string;
  children: ReactNode;
  icon?: LucideIcon;
  active?: boolean;
  className?: string;
}

export default function SidebarSection({
  title,
  children,
  icon: Icon,
  active = false,
  className,
}: SidebarSectionProps) {
  return (
    <section
      className={clsx(
        "mb-8 animate-in fade-in slide-in-from-left-2 duration-300",
        className,
      )}
    >
      <div
        className={clsx(
          `
          mb-3
          flex
          items-center
          gap-2
          px-4
          text-[11px]
          font-semibold
          uppercase
          tracking-[0.22em]
          transition-colors
          `,
          active ? "text-blue-400" : "text-zinc-500",
        )}
      >
        {Icon && (
          <Icon
            className="
            h-4
            w-4
            "
          />
        )}

        <span>{title}</span>
      </div>

      <div
        className="
        space-y-1
        "
      >
        {children}
      </div>
    </section>
  );
}
