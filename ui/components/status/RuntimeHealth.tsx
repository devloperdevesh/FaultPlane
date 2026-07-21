"use client";

const metrics = [
  ["CPU", "42%"],

  ["Memory", "61%"],

  ["Network", "30%"],

  ["Goroutines", "128"],

  ["Sockets", "340"],
];

export default function RuntimeHealth() {
  return (
    <div
      className="
grid
gap-4
md:grid-cols-5
"
    >
      {metrics.map(([name, value]) => (
        <div
          key={name}
          className="
rounded-xl
border
border-white/10
bg-zinc-950
p-5
"
        >
          <p
            className="
text-xs
text-zinc-500
"
          >
            {name}
          </p>

          <p
            className="
mt-3
font-mono
text-xl
text-white
"
          >
            {value}
          </p>
        </div>
      ))}
    </div>
  );
}
