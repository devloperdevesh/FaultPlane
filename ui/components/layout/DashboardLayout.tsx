"use client";

interface DashboardLayoutProps {
  children: React.ReactNode;
  inspector?: React.ReactNode;
}

export default function DashboardLayout({
  children,
  inspector,
}: DashboardLayoutProps) {
  return (
    <div
      className="
        grid
        h-full
        grid-cols-[1fr_320px]
        gap-6
      "
    >
      {/* Main Workspace */}
      <section>{children}</section>

      {/* Right Inspector Panel */}
      <aside
        className="
          rounded-xl
          border
          border-zinc-800
          bg-zinc-950
          p-5
        "
      >
        {inspector ?? (
          <div className="text-sm text-zinc-500">Runtime Inspector</div>
        )}
      </aside>
    </div>
  );
}
