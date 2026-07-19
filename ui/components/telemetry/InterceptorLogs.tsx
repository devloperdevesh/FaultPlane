"use client";

import { motion } from "framer-motion";

const logs = [
  {
    type: "INTERCEPT",
    source: "Worker-01",
    target: "Gateway",
    status: "captured",
  },

  {
    type: "REDIRECT",
    source: "Gateway",
    target: "Worker-03",
    status: "rerouted",
  },

  {
    type: "BLOCK",
    source: "Unknown",
    target: "Sandbox",
    status: "blocked",
  },
];

export default function InterceptorLogs() {
  return (
    <div
      className="
rounded-2xl
border
border-white/10
bg-zinc-950/80
p-6
backdrop-blur-xl
"
    >
      <h2
        className="
text-lg
font-semibold
text-white
"
      >
        Transport Interceptor
      </h2>

      <p
        className="
mt-1
text-sm
text-zinc-500
"
      >
        Network interception and routing decisions
      </p>

      <div
        className="
mt-6
space-y-3
"
      >
        {logs.map((log, index) => (
          <motion.div
            key={index}
            initial={{
              opacity: 0,
              x: -20,
            }}
            animate={{
              opacity: 1,
              x: 0,
            }}
            transition={{
              delay: index * 0.1,
            }}
            className="
flex
items-center
justify-between
rounded-xl
border
border-white/10
bg-zinc-900/60
p-4
"
          >
            <div>
              <p
                className="
text-sm
font-medium
text-white
"
              >
                {log.type}
              </p>

              <p
                className="
mt-1
text-xs
text-zinc-500
"
              >
                {log.source} → {log.target}
              </p>
            </div>

            <span
              className="
rounded-full
bg-emerald-500/10
px-3
py-1
text-xs
text-emerald-400
"
            >
              {log.status}
            </span>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
