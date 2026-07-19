"use client";

export default function SidebarLogo() {
  return (
    <div
      className="
      flex
      h-20
      items-center
      border-b
      border-zinc-800
      px-6
      "
    >
      <div
        className="
        mr-4
        flex
        h-11
        w-11
        items-center
        justify-center
        rounded-xl
        border
        border-blue-500/20
        bg-blue-500/10
        shadow-lg
        shadow-blue-500/10
        "
      >
        <div className="h-3 w-3 rounded-full bg-blue-500" />
      </div>

      <div>
        <h1 className="text-lg font-semibold tracking-tight text-white">
          FaultPlane
        </h1>

        <p className="text-xs text-zinc-500">AI Runtime Control Plane</p>
      </div>
    </div>
  );
}
