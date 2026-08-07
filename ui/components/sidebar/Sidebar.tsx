"use client";

import { usePathname } from "next/navigation";

import SidebarItem from "./SidebarItem";
import SidebarLogo from "./SidebarLogo";
import SidebarSection from "./SidebarSection";

import { navigation, type NavigationItem } from "@/lib/navigation";

function isLeafItem(
  item: NavigationItem,
): item is NavigationItem & { href: string } {
  return Boolean(item.href);
}

function renderNavigation(items: NavigationItem[], pathname: string) {
  return items.map((item) => {
    if (item.children?.length) {
      const hasActiveChild = item.children.some(
        (child) => child.href === pathname,
      );

      return (
        <SidebarSection
          key={item.title}
          title={item.title}
          icon={item.icon}
          active={hasActiveChild}
        >
          {renderNavigation(item.children, pathname)}
        </SidebarSection>
      );
    }

    if (isLeafItem(item)) {
      return (
        <SidebarItem
          key={item.href}
          title={item.title}
          href={item.href}
          icon={item.icon}
          active={pathname === item.href}
        />
      );
    }

    return null;
  });
}

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
      {/* Brand */}

      <SidebarLogo />

      {/* Navigation */}

      <nav
        className="
        flex-1
        overflow-y-auto
        px-4
        py-6
        "
      >
        {renderNavigation(navigation, pathname)}
      </nav>

      {/* Footer */}

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
          "
        >
          <div
            className="
            flex
            items-center
            justify-between
            "
          >
            <span
              className="
              text-sm
              font-medium
              text-white
              "
            >
              Enterprise
            </span>

            <span
              className="
              rounded-full
              border
              border-emerald-500/30
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

          <p
            className="
            mt-2
            text-xs
            text-zinc-500
            "
          >
            FaultPlane v0.1.0
          </p>
        </div>
      </footer>
    </aside>
  );
}
