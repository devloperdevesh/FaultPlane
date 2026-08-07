"use client";

import { ReactNode } from "react";

interface WorkspaceProps {
  children: ReactNode;
}

export default function Workspace({ children }: WorkspaceProps) {
  return (
    <section
      className="
        flex-1
        h-full
        overflow-auto
        bg-zinc-950
      "
    >
      {children}
    </section>
  );
}
