"use client";

import { usePathname } from "next/navigation";

import SidebarItem from "./SidebarItem";
import SidebarLogo from "./SidebarLogo";
import SidebarSection from "./SidebarSection";

import { navigation } from "@/lib/navigation";

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside
      className="
      flex
      h-screen
      w-72
      shrink-0
      flex-col
      border-r
      border-zinc-800
      bg-zinc-950
      "
    >
      <SidebarLogo />

      <nav
        className="
        flex-1
        overflow-y-auto
        px-4
        py-6
        scrollbar-thin
        scrollbar-thumb-zinc-800
        scrollbar-track-transparent
        "
      >
        <SidebarSection title="Operations">
          {navigation.slice(0, 3).map((item) => (
            <SidebarItem
              key={item.href}
              {...item}
              active={pathname === item.href}
            />
          ))}
        </SidebarSection>

        <SidebarSection title="Runtime">
          {navigation.slice(3, 7).map((item) => (
            <SidebarItem
              key={item.href}
              {...item}
              active={pathname === item.href}
            />
          ))}
        </SidebarSection>

        <SidebarSection title="System">
          {navigation.slice(7).map((item) => (
            <SidebarItem
              key={item.href}
              {...item}
              active={pathname === item.href}
            />
          ))}
        </SidebarSection>
      </nav>

      <footer
        className="
        border-t
        border-zinc-800
        p-5
        "
      >
        <div
          className="
          rounded-xl
          border
          border-zinc-800
          bg-zinc-900/70
          p-4
          backdrop-blur
          "
        >
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-white">Enterprise</span>

            <span
              className="
              rounded-full
              bg-emerald-500/10
              px-2
              py-0.5
              text-[10px]
              font-semibold
              uppercase
              tracking-wide
              text-emerald-400
              "
            >
              Stable
            </span>
          </div>

          <p className="mt-2 text-xs text-zinc-500">Version v0.1.0</p>
        </div>
      </footer>
    </aside>
  );
}
