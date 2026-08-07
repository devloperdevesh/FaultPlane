"use client";

import { motion } from "framer-motion";

const snapshots = ["v40", "v41", "v42", "v43"];

export default function SnapshotTimeline() {
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
        Checkpoint History
      </h3>

      <div
        className="
mt-8
flex
items-center
justify-between
"
      >
        {snapshots.map((snapshot, index) => (
          <div
            key={snapshot}
            className="
flex
items-center
"
          >
            <motion.div
              whileHover={{
                scale: 1.1,
              }}
              className="
flex
h-12
w-12
items-center
justify-center
rounded-full
border
border-blue-500/30
bg-blue-500/10
text-xs
text-blue-400
"
            >
              {snapshot}
            </motion.div>

            {index !== snapshots.length - 1 && (
              <div
                className="
mx-3
h-px
w-14
bg-zinc-700
"
              />
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
