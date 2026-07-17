"use client";

import { usePathname } from "next/navigation";

import SidebarItem from "./SidebarItem";
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
      flex-col
      border-r
      border-zinc-800
      bg-zinc-950
      "
    >
      {/* Logo */}

      <div
        className="
        flex
        h-20
        items-center
        border-b
        border-zinc-800
        px-6
        "
      >
        <div
          className="
          mr-3
          flex
          h-10
          w-10
          items-center
          justify-center
          rounded-xl
          bg-blue-500/10
          "
        >
          <div className="h-3 w-3 rounded-full bg-blue-500" />
        </div>

        <div>
          <h1 className="font-semibold text-white">FaultPlane</h1>

          <p className="text-xs text-zinc-500">AI Runtime Control Plane</p>
        </div>
      </div>

      {/* Navigation */}

      <div
        className="
        flex-1
        overflow-y-auto
        px-4
        py-6
        "
      >
        <SidebarSection title="Operations">
          {navigation.slice(0, 3).map((item) => (
            <SidebarItem
              key={item.href}
              title={item.title}
              href={item.href}
              icon={item.icon}
              active={pathname === item.href}
            />
          ))}
        </SidebarSection>

        <SidebarSection title="Runtime">
          {navigation.slice(3, 7).map((item) => (
            <SidebarItem
              key={item.href}
              title={item.title}
              href={item.href}
              icon={item.icon}
              active={pathname === item.href}
            />
          ))}
        </SidebarSection>

        <SidebarSection title="System">
          {navigation.slice(7).map((item) => (
            <SidebarItem
              key={item.href}
              title={item.title}
              href={item.href}
              icon={item.icon}
              active={pathname === item.href}
            />
          ))}
        </SidebarSection>
      </div>

      {/* Footer */}

      <div
        className="
        border-t
        border-zinc-800
        p-5
        "
      >
        <div className="rounded-xl bg-zinc-900 p-4">
          <p className="text-sm font-medium text-white">Enterprise Mode</p>

          <p className="mt-1 text-xs text-zinc-500">Version v0.1.0</p>
        </div>
      </div>
    </aside>
  );
}
