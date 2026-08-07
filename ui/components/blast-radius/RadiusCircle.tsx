"use client";

import { motion } from "framer-motion";

export default function RadiusCircle() {
  return (
    <div
      className="
relative
h-72
w-72
flex
items-center
justify-center
"
    >
      {[1, 2, 3].map((i) => (
        <motion.div
          key={i}
          animate={{
            scale: [1, 1.2, 1],

            opacity: [0.4, 0.1, 0.4],
          }}
          transition={{
            repeat: Infinity,

            duration: 2 + i,
          }}
          className="
absolute
rounded-full
border
border-red-500/40
"
          style={{
            height: `${i * 90}px`,
            width: `${i * 90}px`,
          }}
        />
      ))}

      <div
        className="
z-10
rounded-full
bg-red-500/20
border
border-red-500
p-10
"
      >
        Worker Failure
      </div>
    </div>
  );
}
