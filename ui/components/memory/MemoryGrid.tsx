"use client";

import MemoryPage from "./MemoryPage";
import SwapAnimation from "./SwapAnimation";
import MemoryStats from "./MemoryStats";

const pages = [
  "stable",
  "stable",
  "migrating",

  "migrating",
  "restored",
  "restored",

  "stable",
  "migrating",
  "restored",
] as const;

const legend = [
  {
    label: "Stable",
    color: "bg-cyan-400",
  },
  {
    label: "Migrating",
    color: "bg-amber-400",
  },
  {
    label: "Restored",
    color: "bg-emerald-400",
  },
];

export default function MemoryGrid() {
  return (
    <section
      className="
rounded-2xl
border
border-zinc-800
bg-zinc-950
p-6
shadow-xl
shadow-black/20
"
    >
      {/* Header */}

      <div
        className="
flex
items-start
justify-between
mb-8
"
      >
        <div>
          <h2
            className="
text-base
font-semibold
text-white
"
          >
            Runtime Memory Hot Swap
          </h2>

          <p
            className="
mt-1
text-xs
text-zinc-500
"
          >
            Page-level state migration during worker failover
          </p>
        </div>

        <div
          className="
rounded-full
border
border-emerald-500/20
bg-emerald-500/10
px-3
py-1
"
        >
          <span
            className="
text-xs
font-mono
text-emerald-400
"
          >
            HOT SWAP ACTIVE
          </span>
        </div>
      </div>

      {/* Memory Topology */}

      <div
        className="
grid
gap-8
lg:grid-cols-2
"
      >
        {/* Source Memory */}

        <div>
          <p
            className="
mb-3
text-xs
uppercase
tracking-wider
text-zinc-500
"
          >
            Worker A Memory Snapshot
          </p>

          <div
            className="
relative
grid
grid-cols-3
gap-3
rounded-xl
border
border-zinc-800
bg-black/40
p-5
"
          >
            {pages.map((page, index) => (
              <MemoryPage key={index} state={page} />
            ))}

            <SwapAnimation />
          </div>
        </div>

        {/* Migration Info */}

        <div
          className="
flex
flex-col
justify-center
space-y-4
"
        >
          <div
            className="
rounded-xl
border
border-zinc-800
bg-zinc-900/50
p-5
"
          >
            <h3
              className="
text-sm
font-medium
text-white
"
            >
              Page Migration Flow
            </h3>

            <div
              className="
mt-4
space-y-3
text-xs
font-mono
"
            >
              <div className="text-zinc-400">
                Snapshot Captured
                <span className="float-right text-cyan-400">✓</span>
              </div>

              <div className="text-zinc-400">
                Dirty Pages Synced
                <span className="float-right text-amber-400">248</span>
              </div>

              <div className="text-zinc-400">
                Target Restored
                <span className="float-right text-emerald-400">✓</span>
              </div>
            </div>
          </div>

          <div
            className="
flex
gap-4
"
          >
            {legend.map((item) => (
              <div
                key={item.label}
                className="
flex
items-center
gap-2
"
              >
                <div
                  className={`
h-2
w-2
rounded-full
${item.color}
`}
                />

                <span
                  className="
text-xs
text-zinc-500
"
                >
                  {item.label}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Metrics */}

      <div
        className="
mt-8
grid
gap-4
md:grid-cols-3
"
      >
        <MemoryStats title="Snapshot Size" value="1.8 GB" />

        <MemoryStats
          title="Migration Latency"
          value="14.2 ms"
          status="success"
        />

        <MemoryStats title="Runtime State" value="RESTORED" status="success" />
      </div>
    </section>
  );
}
