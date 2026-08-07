"use client";

import { motion } from "framer-motion";

const streams = [
  "Client → Gateway",

  "Gateway → Proxy Layer",

  "Proxy → Worker-02",

  "Worker → Checkpoint Store",
];

export default function ProxyStream() {
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
      <h3
        className="
text-sm
font-semibold
text-white
"
      >
        Proxy Stream
      </h3>

      <div
        className="
mt-6
space-y-4
"
      >
        {streams.map((item, index) => (
          <div
            key={item}
            className="
relative
flex
items-center
gap-4
"
          >
            <div
              className="
h-3
w-3
rounded-full
bg-emerald-400
"
            />

            <p
              className="
text-sm
text-zinc-300
"
            >
              {item}
            </p>

            {index !== streams.length - 1 && (
              <motion.div
                animate={{
                  height: ["0%", "100%"],
                }}
                transition={{
                  duration: 2,
                  repeat: Infinity,
                }}
                className="
absolute
left-[5px]
top-4
w-px
bg-emerald-500/40
"
              />
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
