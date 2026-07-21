"use client";

import { useEffect, useState } from "react";

export default function CommandPalette() {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    function handler(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();

        setOpen(true);
      }
    }

    window.addEventListener("keydown", handler);

    return () => window.removeEventListener("keydown", handler);
  }, []);

  if (!open) return null;

  return (
    <div
      className="
fixed
inset-0
z-50
flex
items-start
justify-center
pt-32
bg-black/60
"
    >
      <div
        className="
w-96
rounded-xl
border
border-white/10
bg-zinc-950
p-4
"
      >
        <p className="text-white">Search FaultPlane</p>

        <div className="mt-4 space-y-2 text-sm text-zinc-400">
          <div>Open Workers</div>

          <div>View eBPF Trace</div>

          <div>Run Chaos Simulation</div>

          <div>Inspect Checkpoint</div>
        </div>
      </div>
    </div>
  );
}
