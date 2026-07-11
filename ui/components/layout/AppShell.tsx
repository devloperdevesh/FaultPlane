"use client";

import { ReactNode } from "react";

type Props = {
  sidebar: ReactNode;
  navbar: ReactNode;
  children: ReactNode;
};

export default function AppShell({ sidebar, navbar, children }: Props) {
  return (
    <div className="flex h-screen bg-[#09090B] text-white">
      <aside className="w-72 border-r border-zinc-800">{sidebar}</aside>

      <main className="flex flex-1 flex-col overflow-hidden">
        <header className="h-16 border-b border-zinc-800">{navbar}</header>

        <section className="flex-1 overflow-y-auto p-6">{children}</section>
      </main>
    </div>
  );
}
