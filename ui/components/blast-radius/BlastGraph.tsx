"use client";

const nodes = [
  {
    name: "Gateway",
    status: "healthy",
  },
  {
    name: "Worker-01",
    status: "failed",
  },
  {
    name: "Worker-02",
    status: "degraded",
  },
  {
    name: "Worker-03",
    status: "healthy",
  },
];

const colors = {
  healthy: "border-emerald-500 bg-emerald-500/10 text-emerald-400",

  degraded: "border-orange-500 bg-orange-500/10 text-orange-400",

  failed: "border-red-500 bg-red-500/10 text-red-400",
};

export default function BlastGraph() {
  return (
    <div
      className="
rounded-xl
border
border-zinc-800
bg-black/40
p-6
"
    >
      <h3
        className="
mb-6
text-sm
font-semibold
text-white
"
      >
        Failure Dependency Graph
      </h3>

      <div
        className="
flex
flex-col
items-center
gap-6
"
      >
        {nodes.map((node) => (
          <div
            key={node.name}
            className={`
w-40
rounded-xl
border
px-4
py-3
text-center
transition
hover:scale-105
${colors[node.status as keyof typeof colors]}
`}
          >
            <p
              className="
text-sm
font-medium
"
            >
              {node.name}
            </p>

            <p
              className="
mt-1
text-xs
capitalize
"
            >
              {node.status}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
