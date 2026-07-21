"use client";

import { motion } from "framer-motion";

import WorkerImpact from "./WorkerImpact";
import BlastGraph from "./BlastGraph";
import ImpactTimeline from "./ImpactTimeline";
import RadiusCircle from "./RadiusCircle";

const workers = [
  {
    name: "Gateway",
    status: "healthy",
    layer: "entry",
  },

  {
    name: "Worker-01",
    status: "failed",
    layer: "compute",
  },

  {
    name: "Worker-02",
    status: "degraded",
    layer: "compute",
  },

  {
    name: "Checkpoint Store",
    status: "recovering",
    layer: "storage",
  },
] as const;

const metrics = [
  {
    title: "Affected Nodes",
    value: "2 / 24",
    status: "critical",
  },

  {
    title: "Recovery Progress",
    value: "87%",
    status: "healthy",
  },

  {
    title: "Containment",
    value: "ACTIVE",
    status: "healthy",
  },

  {
    title: "MTTR",
    value: "14.2ms",
    status: "healthy",
  },
];

export default function BlastRadius() {
  return (
    <section
      className="
relative
overflow-hidden
rounded-2xl
border
border-white/10
bg-zinc-950/80
backdrop-blur-xl
p-8
space-y-8
"
    >
      {/* Header */}

      <div
        className="
flex
items-start
justify-between
"
      >
        <div>
          <h2
            className="
text-lg
font-semibold
text-white
"
          >
            Recovery Blast Radius
          </h2>

          <p
            className="
mt-1
text-sm
text-zinc-500
"
          >
            Failure propagation, impact analysis and recovery containment
          </p>
        </div>

        <div
          className="
rounded-full
border
border-emerald-500/30
bg-emerald-500/10
px-4
py-2
text-xs
text-emerald-400
"
        >
          Recovery Active
        </div>
      </div>

      {/* Visualization */}

      <div
        className="
grid
gap-8
xl:grid-cols-3
"
      >
        {/* Radius */}

        <div
          className="
flex
items-center
justify-center
rounded-xl
border
border-white/10
bg-black/40
p-8
"
        >
          <RadiusCircle />
        </div>

        {/* Dependency Graph */}

        <div
          className="
rounded-xl
border
border-white/10
bg-black/40
p-6
"
        >
          <h3
            className="
mb-5
text-sm
font-medium
text-white
"
          >
            Failure Dependency Graph
          </h3>

          <BlastGraph />
        </div>

        {/* Timeline */}

        <div
          className="
rounded-xl
border
border-white/10
bg-black/40
p-6
"
        >
          <h3
            className="
mb-5
text-sm
font-medium
text-white
"
          >
            Impact Timeline
          </h3>

          <ImpactTimeline />
        </div>
      </div>

      {/* Worker Impact */}

      <div
        className="
grid
gap-4
md:grid-cols-2
xl:grid-cols-4
"
      >
        {workers.map((worker, index) => (
          <motion.div
            key={worker.name}
            initial={{
              opacity: 0,
              y: 20,
            }}
            animate={{
              opacity: 1,
              y: 0,
            }}
            transition={{
              delay: index * 0.08,
            }}
          >
            <WorkerImpact name={worker.name} status={worker.status} />
          </motion.div>
        ))}
      </div>

      {/* Metrics */}

      <div
        className="
grid
gap-4
md:grid-cols-2
xl:grid-cols-4
"
      >
        {metrics.map((metric) => (
          <div
            key={metric.title}
            className="
rounded-xl
border
border-white/10
bg-zinc-900/70
p-5
"
          >
            <p
              className="
text-xs
uppercase
tracking-wider
text-zinc-500
"
            >
              {metric.title}
            </p>

            <p
              className={`
mt-3
font-mono
text-xl

${metric.status === "critical" ? "text-red-400" : "text-emerald-400"}

`}
            >
              {metric.value}
            </p>
          </div>
        ))}
      </div>

      {/* Propagation overlay */}

      <motion.div
        animate={{
          opacity: [0.2, 0.5, 0.2],
        }}
        transition={{
          duration: 3,

          repeat: Infinity,
        }}
        className="
pointer-events-none
absolute
inset-0
bg-gradient-to-r
from-red-500/5
via-transparent
to-emerald-500/5
"
      />
    </section>
  );
}
