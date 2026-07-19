"use client";

const events = [
  "syscall hook attached",

  "socket migration detected",

  "memory page swapped",

  "kernel trace started",
];

export default function KernelLogs() {
  return (
    <div
      className="
rounded-2xl
border
border-white/10
bg-zinc-950
p-6
"
    >
      <h2
        className="
mb-5
text-sm
font-semibold
text-white
"
      >
        Kernel Events
      </h2>

      <div
        className="
space-y-3
"
      >
        {events.map((event) => (
          <div
            key={event}
            className="
rounded-lg
bg-zinc-900
p-3
font-mono
text-xs
text-zinc-300
"
          >
            {k > ">"}
            {event}
          </div>
        ))}
      </div>
    </div>
  );
}
