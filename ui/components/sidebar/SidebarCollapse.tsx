"use client";

import { ChevronLeft, ChevronRight } from "lucide-react";

interface Props {
  collapsed: boolean;
  onToggle: () => void;
}

export default function SidebarCollapse({
  collapsed,
  onToggle,
}: Props) {
  const Icon = collapsed ? ChevronRight : ChevronLeft;

  return (
    <button
      type="button"
      onClick={onToggle}
      aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      className="
        absolute
        -right-4
        top-[68px]
        z-50
        flex
        h-8
        w-8
        items-center
        justify-center
        rounded-full
        border
        border-emerald-500/40
        bg-zinc-950
        text-zinc-400
        shadow-xl
        transition-all
        duration-200
        hover:border-emerald-400
        hover:bg-zinc-900
        hover:text-white
        focus:outline-none
        focus:ring-2
        focus:ring-emerald-400/40
      "
    >
      <Icon size={16} strokeWidth={2} />
    </button>
  );
}
