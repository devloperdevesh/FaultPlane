"use client";

import { usePathname } from "next/navigation";
import { useState } from "react";
import type { ReactNode } from "react";

import SidebarItem from "./SidebarItem";
import SidebarLogo from "./SidebarLogo";
import SidebarSection from "./SidebarSection";
import SidebarCollapse from "./SidebarCollapse";

import { navigation, type NavigationItem } from "@/lib/navigation";

function isLeafItem(
  item: NavigationItem,
): item is NavigationItem & { href: string } {
  return Boolean(item.href);
}

function renderNavigation(
  items: NavigationItem[],
  pathname: string,
  collapsed: boolean,
): ReactNode[] {
  return items.flatMap((item): ReactNode[] => {
    if (item.children?.length) {
      const hasActiveChild = item.children.some(
        (child) => child.href === pathname,
      );

      if (collapsed) {
        return renderNavigation(
          item.children,
          pathname,
          true,
        );
      }

      return [
        <SidebarSection
          key={item.title}
          title={item.title}
          icon={item.icon}
          active={hasActiveChild}
        >
          {renderNavigation(
            item.children,
            pathname,
            false,
          )}
        </SidebarSection>,
      ];
    }

    if (isLeafItem(item)) {
      return [
        <SidebarItem
          key={item.href}
          title={item.title}
          href={item.href}
          icon={item.icon}
          active={pathname === item.href}
          collapsed={collapsed}
        />,
      ];
    }

    return [];
  });
}

export default function Sidebar() {
  const pathname = usePathname();
  const [collapsed, setCollapsed] = useState(false);

  return (
    <aside
      className={`relative flex h-screen shrink-0 flex-col border-r border-zinc-800 bg-zinc-950 transition-[width] duration-200 ${
        collapsed ? "w-20" : "w-72"
      }`}
    >
      <SidebarLogo collapsed={collapsed} />

      <SidebarCollapse
        collapsed={collapsed}
        onToggle={() => setCollapsed((value) => !value)}
      />

      <nav
        className={`flex-1 overflow-y-auto py-6 ${
          collapsed ? "px-2" : "px-4"
        }`}
      >
        <div className="space-y-1">
          {renderNavigation(
            navigation,
            pathname,
            collapsed,
          )}
        </div>
      </nav>

      {!collapsed && (
        <footer className="border-t border-zinc-800 p-5">
          <div className="rounded-xl border border-zinc-800 bg-zinc-900/70 p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-white">
                Enterprise
              </span>

              <span className="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-emerald-400">
                Stable
              </span>
            </div>

            <p className="mt-2 text-xs text-zinc-500">
              FaultPlane v0.1.0
            </p>
          </div>
        </footer>
      )}
    </aside>
  );
}
