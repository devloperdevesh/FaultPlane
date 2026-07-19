"use client";

import { motion } from "framer-motion";

const logs = [
  {
    level: "INFO",
    message: "checkpoint created",
    time: "10:01:02",
  },

  {
    level: "SUCCESS",
    message: "recovery completed",
    time: "10:01:14",
  },

  {
    level: "WARN",
    message: "worker degraded",
    time: "10:01:20",
  },

  {
    level: "ERROR",
    message: "connection timeout",
    time: "10:01:32",
  },
];

const colors = {
  INFO: "text-blue-400",
  SUCCESS: "text-emerald-400",
  WARN: "text-yellow-400",
  ERROR: "text-red-400",
};

export default function TelemetryLogs() {
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
      <div
        className="
mb-5
flex
justify-between
"
      >
        <h2
          className="
text-sm
font-semibold
text-white
"
        >
          Runtime Telemetry
        </h2>

        <span
          className="
text-xs
text-emerald-400
"
        >
          LIVE
        </span>
      </div>

      <div className="space-y-3">
        {logs.map((log, index) => (
          <motion.div
            key={index}
            initial={{
              opacity: 0,
              x: -10,
            }}
            animate={{
              opacity: 1,
              x: 0,
            }}
            transition={{
              delay: index * 0.05,
            }}
            className="
flex
items-center
justify-between
rounded-xl
bg-zinc-900
p-4
font-mono
text-sm
"
          >
            <div>
              <span className={colors[log.level as keyof typeof colors]}>
                [{log.level}]
              </span>

              <span
                className="
ml-3
text-zinc-300
"
              >
                {log.message}
              </span>
            </div>

            <span
              className="
text-xs
text-zinc-500
"
            >
              {log.time}
            </span>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
