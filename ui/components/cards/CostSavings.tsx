"use client";

import { motion } from "framer-motion";

const data = [
  {
    label: "Monthly Spend",
    value: "$4,820",
  },
  {
    label: "Potential Saving",
    value: "$1,240",
  },
  {
    label: "Optimization",
    value: "32%",
  },
];

export default function CostSavings() {
  return (
    <section
      className="
rounded-2xl
border
border-white/10
bg-zinc-950/80
backdrop-blur-xl
p-6
"
    >
      <div className="mb-6">
        <h2
          className="
text-lg
font-semibold
text-white
"
        >
          FinOps Intelligence
        </h2>

        <p
          className="
mt-1
text-sm
text-zinc-500
"
        >
          Runtime cloud cost optimization
        </p>
      </div>

      <div
        className="
grid
gap-4
md:grid-cols-3
"
      >
        {data.map((item, index) => (
          <motion.div
            key={item.label}
            initial={{
              opacity: 0,
              y: 20,
            }}
            animate={{
              opacity: 1,
              y: 0,
            }}
            transition={{
              delay: index * 0.1,
            }}
            whileHover={{
              scale: 1.03,
            }}
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
tracking-wide
text-zinc-500
"
            >
              {item.label}
            </p>

            <p
              className="
mt-3
font-mono
text-2xl
text-emerald-400
"
            >
              {item.value}
            </p>
          </motion.div>
        ))}
      </div>
    </section>
  );
}
