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
  FlaskConical,
  Clock3,
  Shield,
  Sparkles,
  DollarSign,
} from "lucide-react";


export interface NavigationItem {
  title: string;
  href?: string;
  icon: LucideIcon;

  children?: NavigationItem[];
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
    title: "Experimental",
    icon: FlaskConical,

    children: [

      {
        title: "Time Travel",
        href: "/experimental/time-travel",
        icon: Clock3,
      },


      {
        title: "WASM Sandbox",
        href: "/experimental/wasm",
        icon: Shield,
      },


      {
        title: "Speculative Paths",
        href: "/experimental/speculative",
        icon: Sparkles,
      },


      {
        title: "FinOps",
        href: "/experimental/finops",
        icon: DollarSign,
      },

    ],
  },


  {
    title: "Settings",
    href: "/settings",
    icon: Settings,
  },

];