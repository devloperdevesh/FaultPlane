"use client";

import "./globals.css";
import React, { useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";

import {
  LayoutDashboard,
  Users,
  Network,
  AlertTriangle,
  Database,
  RefreshCw,
  Layers,
  ShieldCheck,
  ChevronLeft,
  ChevronRight,
  Activity,
  Search,
  Bell,
  ChevronDown,
} from "lucide-react";


export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [collapsed, setCollapsed] = useState(false);

  const pathname = usePathname();

  const menu = [
    {
      name: "Operations Console",
      icon: LayoutDashboard,
      path: "/dashboard",
    },

    {
      name: "Agent Fleet",
      icon: Users,
      path: "/dashboard/agents",
    },

    {
      name: "Workflow Control",
      icon: Network,
      path: "/dashboard/workflows",
    },

    {
      name: "Incident Center",
      icon: AlertTriangle,
      path: "/dashboard/incidents",
      badge: "2",
    },

    {
      name: "Checkpoint Store",
      icon: Database,
      path: "/dashboard/checkpoints",
    },

    {
      name: "Recovery Center",
      icon: RefreshCw,
      path: "/dashboard/recovery",
    },

    {
      name: "Execution Explorer",
      icon: Layers,
      path: "/dashboard/traces",
    },

    {
      name: "Telemetry",
      icon: ShieldCheck,
      path: "/dashboard/telemetry",
    },
  ];

  return (
    <html lang="en">
      <body
        className="
h-screen
overflow-hidden

bg-[#060914]

text-slate-200

antialiased
"
      >
        <div className="flex h-screen">
          {/* SIDEBAR */}

          <aside
            className={`
relative

flex
flex-col
justify-between


${collapsed ? "w-[72px]" : "w-[250px]"}


bg-[#090d18]

border-r

border-white/[0.06]


transition-[width]

duration-300

ease-out


shadow-xl

z-40

`}
          >
            <div>
              {/* BRAND */}

              <div
                className="
h-[76px]

px-4

flex

items-center

gap-3


border-b

border-white/[0.06]

"
              >
                <div
                  className="
w-10
h-10

rounded-xl

bg-white

flex

items-center

justify-center

overflow-hidden

shrink-0

shadow-sm

"
                >
                  <Image
                    src="/favicon.ico"
                    alt="FaultPlane"
                    width={34}
                    height={34}
                    className="
object-contain
scale-125
"
                  />
                </div>

                {!collapsed && (
                  <div>
                    <h1
                      className="
text-[17px]

font-semibold

tracking-tight

text-white

"
                    >
                      FaultPlane
                    </h1>

                    <p
                      className="
text-[9px]

uppercase

tracking-[0.25em]

text-slate-500

"
                    >
                      AI RELIABILITY PLATFORM
                    </p>
                  </div>
                )}
              </div>

              <nav
                className="
px-3

py-5

space-y-1

"
              >
                {menu.map((item) => {
                  const Icon = item.icon;

                  const active = pathname === item.path;

                  return (
                    <Link key={item.path} href={item.path}>
                      <div
                        className={`

relative

flex

items-center


${collapsed ? "justify-center" : "justify-between"}


px-3

py-3


rounded-lg


transition-all

duration-200


group


cursor-pointer


${
  active
    ? "bg-blue-500/[0.08] text-blue-400"
    : "text-slate-400 hover:text-white hover:bg-white/[0.04]"
}

`}
                      >
                        {active && (
                          <div
                            className="
absolute

left-0

h-5

w-[3px]

rounded-r

bg-blue-500

"
                          />
                        )}

                        <div
                          className="
flex

items-center

gap-3

"
                        >
                          <Icon
                            size={17}
                            strokeWidth={1.8}
                            className={
                              active
                                ? "text-blue-400"
                                : "text-slate-500 group-hover:text-slate-300"
                            }
                          />

                          {!collapsed && (
                            <span
                              className="
text-sm

font-medium

"
                            >
                              {item.name}
                            </span>
                          )}
                        </div>

                        {item.badge && !collapsed && (
                          <span
                            className="
text-[11px]

px-2

py-0.5

rounded-full

bg-red-500/10

text-red-400

border

border-red-500/20

font-mono

"
                          >
                            {item.badge}
                          </span>
                        )}
                      </div>
                    </Link>
                  );
                })}
              </nav>
            </div>
            {/* SYSTEM STATUS FOOTER */}

            <div
              className="
p-4

border-t

border-white/[0.06]

"
            >
              <div
                className="
rounded-xl

bg-[#0c111f]

border

border-white/[0.06]

p-3

"
              >
                <div
                  className="
flex

items-center

gap-3

"
                >
                  <span
                    className="
relative

flex

h-2.5

w-2.5

"
                  >
                    <span
                      className="
absolute

inset-0

rounded-full

bg-emerald-400

opacity-40

animate-ping

"
                    />

                    <span
                      className="
relative

h-2.5

w-2.5

rounded-full

bg-emerald-500

"
                    />
                  </span>

                  {!collapsed && (
                    <div>
                      <p
                        className="
text-xs

font-semibold

text-emerald-400

"
                      >
                        SYSTEM OPERATIONAL
                      </p>

                      <p
                        className="
mt-1

text-[11px]

text-slate-500

font-mono

"
                      >
                        Runtime v0.1.0
                      </p>
                    </div>
                  )}
                </div>

                {!collapsed && (
                  <div
                    className="
mt-3

flex

items-center

gap-2

text-[11px]

text-slate-500

font-mono

"
                  >
                    <Activity size={12} />
                    99.98% uptime
                  </div>
                )}
              </div>
            </div>

            {/* COLLAPSE BUTTON */}

            <button
              onClick={() => setCollapsed(!collapsed)}
              className="

absolute

right-[-14px]

top-20


w-7

h-7


rounded-full


bg-[#111827]


border

border-white/10


flex

items-center

justify-center


text-slate-400


hover:text-white


hover:bg-blue-500/20


hover:scale-110


active:scale-95


transition-all

duration-200


shadow-lg

"
            >
              {collapsed ? (
                <ChevronRight size={15} />
              ) : (
                <ChevronLeft size={15} />
              )}
            </button>
          </aside>

          {/* MAIN AREA */}

          <main
            className="

flex-1

h-screen

overflow-hidden

flex

flex-col

bg-[#060914]

"
          >
            {/* TOP BAR */}

            <header
              className="

h-16

shrink-0


border-b

border-white/[0.06]


bg-[#060914]/90


backdrop-blur-xl


flex

items-center

justify-between


px-6


"
            >
              {/* LEFT SIDE */}

              <div
                className="
flex

items-center

gap-4

"
              >
                <h2
                  className="
text-sm

font-semibold

text-white

"
                >
                  Operations Console
                </h2>

                <div
                  className="
hidden

md:flex

items-center

gap-2

px-3

py-1.5


rounded-lg


bg-white/[0.03]


border

border-white/[0.08]


text-xs

text-slate-400

"
                >
                  Production
                  <ChevronDown size={13} />
                </div>
              </div>

              {/* RIGHT SIDE */}

              <div
                className="
flex

items-center

gap-3

"
              >
                <button
                  className="

h-9

px-3


rounded-lg


bg-white/[0.03]


border

border-white/[0.08]


flex

items-center

gap-2


text-xs


text-slate-400


hover:bg-white/[0.06]


transition

"
                >
                  <Search size={14} />

                  <span>Search</span>

                  <span
                    className="
text-[10px]

text-slate-600

font-mono

"
                  >
                    ⌘K
                  </span>
                </button>

                <button
                  className="

relative

w-9

h-9


rounded-lg


border

border-white/[0.08]


flex

items-center

justify-center


text-slate-400


hover:text-white


hover:bg-white/[0.05]


transition


"
                >
                  <Bell size={15} />

                  <span
                    className="

absolute

top-2

right-2


w-2

h-2


rounded-full


bg-red-500

"
                  />
                </button>
              </div>
            </header>

            {/* PAGE CONTENT */}

            <div
              className="

flex-1

overflow-y-auto


px-6

py-6


"
            >
              <div
                className="

max-w-[1600px]

mx-auto

"
              >
                {children}
              </div>
            </div>
          </main>
        </div>
      </body>
    </html>
  );
}