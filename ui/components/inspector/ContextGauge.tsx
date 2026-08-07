"use client";

interface Props {
  value: number;
}

export default function ContextGauge({ value }: Props) {
  return (
    <div className="space-y-2">
      <div className="flex justify-between text-xs text-zinc-400">
        <span>Context Usage</span>

        <span>{value}%</span>
      </div>

      <div
        className="
h-2
rounded-full
bg-zinc-800
overflow-hidden
"
      >
        <div
          style={{
            width: `${value}%`,
          }}
          className="
h-full
rounded-full
bg-blue-500
"
        />
      </div>
    </div>
  );
}
