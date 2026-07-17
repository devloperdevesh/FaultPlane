import type { LucideIcon } from "lucide-react";

import {
  Activity,
  BarChart3,
  Database,
  LayoutDashboard,
  RefreshCcw,
  Server,
  Settings,
  Workflow,
} from "lucide-react";

export interface NavigationItem {
  title: string;
  href: string;
  icon: LucideIcon;
}

export const navigation: NavigationItem[] = [
  {
    title: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    title: "Workers",
    href: "/workers",
    icon: Server,
  },
  {
    title: "Workflows",
    href: "/workflows",
    icon: Workflow,
  },
  {
    title: "Checkpoints",
    href: "/checkpoints",
    icon: Database,
  },
  {
    title: "Recovery",
    href: "/recovery",
    icon: RefreshCcw,
  },
  {
    title: "Telemetry",
    href: "/telemetry",
    icon: Activity,
  },
  {
    title: "Metrics",
    href: "/metrics",
    icon: BarChart3,
  },
  {
    title: "Settings",
    href: "/settings",
    icon: Settings,
  },
];