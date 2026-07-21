"use client";

import {
  Activity,
  Shield,
  Network,
  Database,
  DollarSign,
  Cpu,
  GitBranch,
} from "lucide-react";

import SidebarItem from "../sidebar/SidebarItem";

const items = [
  {
    title: "Overview",
    href: "/dashboard",
    icon: Activity,
  },

  {
    title: "Recovery",
    href: "/dashboard/recovery",
    icon: Shield,
  },

  {
    title: "Topology",
    href: "/dashboard/topology",
    icon: Network,
  },

  {
    title: "Memory",
    href: "/dashboard/memory",
    icon: Database,
  },

  {
    title: "FinOps",
    href: "/dashboard/finops",
    icon: DollarSign,
  },

  {
    title: "eBPF",
    href: "/dashboard/ebpf",
    icon: Cpu,
  },

  {
    title: "Workflows",
    href: "/dashboard/workflows",
    icon: GitBranch,
  },
];

export default function DashboardSidebar() {
  return (
    <aside
      className="
w-60
border-r
border-white/10
bg-black/40
p-4
"
    >
      <p
        className="
mb-5
px-3
text-xs
uppercase
tracking-widest
text-zinc-500
"
      >
        Dashboard
      </p>

      <div className="space-y-1">
        {items.map((item) => (
          <SidebarItem key={item.href} {...item} />
        ))}
      </div>
    </aside>
  );
}
