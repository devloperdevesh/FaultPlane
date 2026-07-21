"use client";

const checkpoints = [
  {
    id: "cp-1024",
    status: "active",
    size: "1.8GB",
    worker: "worker-01",
  },
  {
    id: "cp-1023",
    status: "expired",
    size: "1.6GB",
    worker: "worker-02",
  },
];

export default function CheckpointDiff() {
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
        Checkpoint Comparison
      </h2>

      <div
        className="
mt-5
space-y-3
"
      >
        {checkpoints.map((cp) => (
          <div
            key={cp.id}
            className="
rounded-xl
border
border-white/10
bg-zinc-900
p-4
"
          >
            <div
              className="
flex
justify-between
"
            >
              <span className="text-white">{cp.id}</span>

              <span className="text-emerald-400 text-xs">{cp.status}</span>
            </div>

            <div
              className="
mt-3
grid
grid-cols-2
text-xs
text-zinc-500
"
            >
              <p>Worker</p>

              <p>Size</p>

              <p className="text-white">{cp.worker}</p>

              <p className="text-white">{cp.size}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
