"use client";

import Link from "next/link";
import clsx from "clsx";
import type { LucideIcon } from "lucide-react";

interface Props {
  title: string;
  href: string;
  icon: LucideIcon;
  active?: boolean;
  badge?: string;
}

export default function SidebarItem({
  title,
  href,
  icon: Icon,
  active,
  badge,
}: Props) {
  return (
    <Link
      href={href}
      className={clsx(
        "group relative flex items-center gap-3 rounded-xl px-4 py-3 text-sm transition-all duration-200",

        active
          ? "bg-white/10 text-white shadow-lg"
          : "text-zinc-400 hover:bg-white/5 hover:text-white",
      )}
    >
      {active && (
        <span
          className="
absolute
left-0
h-7
w-1
rounded-r-full
bg-emerald-400
"
        />
      )}

      <Icon
        size={18}
        className={clsx(
          "transition-transform duration-200 group-hover:scale-110",

          active ? "text-emerald-400" : "text-zinc-500",
        )}
      />

      <span className="flex-1">{title}</span>

      {badge && (
        <span
          className="
rounded-md
bg-zinc-800
px-2
py-1
text-[10px]
"
        >
          {badge}
        </span>
      )}
    </Link>
  );
}
