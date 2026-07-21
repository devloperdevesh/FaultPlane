"use client";

import { ReactNode } from "react";

import Workspace from "./Workspace";
import InspectorPanel from "./InspectorPanel";
import BottomStatusBar from "./BottomStatusBar";

interface AppShellProps {
  children: ReactNode;
}

export default function AppShell({ children }: AppShellProps) {
  return (
    <main
      className="
        h-screen
        w-full
        flex
        flex-col
        overflow-hidden
        bg-zinc-950
        text-zinc-100
      "
    >
      {/* Main Application Area */}

      <section
        className="
          flex-1
          min-h-0
          flex
          overflow-hidden
        "
      >
        {/* Workspace / Canvas */}

        <div
          className="
            flex-1
            min-w-0
            overflow-hidden
          "
        >
          <Workspace>{children}</Workspace>
        </div>

        {/* Runtime Inspector */}

        <aside
          className="
            w-80
            border-l
            border-zinc-800
            bg-zinc-900/40
            overflow-y-auto
          "
        >
          <InspectorPanel />
        </aside>
      </section>

      {/* System Status */}

      <footer
        className="
          h-8
          shrink-0
          border-t
          border-zinc-800
          bg-zinc-950
        "
      >
        <BottomStatusBar />
      </footer>
    </main>
  );
}
