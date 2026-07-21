"use client";

const snapshot = {
  memory: "4.2GB",
  agents: "24",
  tokens: "5100",
  version: "v1.4.2",
};

export default function SnapshotViewer() {
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
text-sm
font-semibold
text-white
"
      >
        Runtime Snapshot
      </h2>

      <div
        className="
mt-5
grid
grid-cols-2
gap-4
"
      >
        {Object.entries(snapshot).map(([key, value]) => (
          <div
            key={key}
            className="
rounded-xl
bg-zinc-900
p-4
"
          >
            <p
              className="
text-xs
text-zinc-500
uppercase
"
            >
              {key}
            </p>

            <p
              className="
mt-2
font-mono
text-white
"
            >
              {value}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
