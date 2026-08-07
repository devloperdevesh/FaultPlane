"use client";

import SimulationResult from "./SimulationResult";

interface Props {
  onClose: () => void;
}

const actions = [
  "Kill Worker",
  "Network Drop",
  "Memory Leak",
  "SIGTERM",
  "SIGKILL",
];

export default function ChaosModal({ onClose }: Props) {
  return (
    <div
      className="
fixed
inset-0
flex
items-center
justify-center
bg-black/70
"
    >
      <div
        className="
w-full
max-w-md
rounded-xl
border
border-zinc-800
bg-zinc-950
p-6
"
      >
        <div className="flex justify-between mb-6">
          <h2 className="text-white font-semibold">Chaos Simulator</h2>

          <button onClick={onClose} className="text-zinc-400">
            X
          </button>
        </div>

        <div className="space-y-3">
          {actions.map((action) => (
            <button
              key={action}
              className="
w-full
rounded-lg
bg-zinc-900
px-4
py-3
text-left
text-sm
text-white
transition
hover:bg-zinc-800
"
            >
              {action}
            </button>
          ))}
        </div>

        <SimulationResult />
      </div>
    </div>
  );
}
