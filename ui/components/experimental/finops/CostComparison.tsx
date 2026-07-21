"use client";

import { BarChart, Bar, XAxis, YAxis, ResponsiveContainer } from "recharts";

const data = [
  {
    name: "Runtime",
    before: 420,
    after: 85,
  },
];

export default function CostComparison() {
  return (
    <div
      className="
h-64
rounded-xl
border
border-white/10
bg-zinc-950
p-5
"
    >
      <ResponsiveContainer>
        <BarChart data={data}>
          <XAxis dataKey="name" />

          <YAxis />

          <Bar dataKey="before" />

          <Bar dataKey="after" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
