"use client";

import { motion } from "framer-motion";

const events = [
  {
    message: "syscall hook attached",
    time: "10:01:02",
  },

  {
    message: "socket migration detected",
    time: "10:01:14",
  },

  {
    message: "memory page swapped",
    time: "10:01:20",
  },

  {
    message: "kernel trace started",
    time: "10:01:32",
  },
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
      <div
        className="
        mb-5
        flex
        items-center
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
          Kernel Events
        </h2>

        <span
          className="
          rounded-full
          bg-blue-500/10
          px-3
          py-1
          text-xs
          text-blue-400
          "
        >
          EXPERIMENTAL
        </span>
      </div>

      <div
        className="
        space-y-3
        "
      >
        {events.map((event, index) => (
          <motion.div
            key={event.message}
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
              text-xs
              "
          >
            <div>
              <span
                className="
                  mr-3
                  text-blue-400
                  "
              >
                {">"}
              </span>

              <span
                className="
                  text-zinc-300
                  "
              >
                {event.message}
              </span>
            </div>

            <span
              className="
                text-zinc-500
                "
            >
              {event.time}
            </span>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
