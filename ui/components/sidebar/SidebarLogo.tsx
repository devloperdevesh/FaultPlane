"use client";

interface Props {
  collapsed?: boolean;
}

export default function SidebarLogo({ collapsed = false }: Props) {
  return (
    <div
      className={`flex h-20 shrink-0 items-center border-b border-white/[0.08] ${
        collapsed ? "justify-center px-3" : "gap-3 px-5"
      }`}
    >
      <div className="relative shrink-0">
        <img
          src="/favicon.ico"
          alt="FaultPlane"
          width={42}
          height={42}
          className="h-[42px] w-[42px] rounded-xl border border-white/[0.08] object-cover shadow-[0_0_30px_rgba(255,255,255,0.06)]"
        />
      </div>

      {!collapsed && (
        <div className="min-w-0">
          <div className="truncate text-[15px] font-semibold tracking-[-0.02em] text-white">
            FaultPlane
          </div>

          <div className="mt-0.5 truncate text-[11px] tracking-wide text-zinc-500">
            AI Runtime Control Plane
          </div>
        </div>
      )}
    </div>
  );
}
