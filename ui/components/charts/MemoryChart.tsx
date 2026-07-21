"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const memoryData = [
  {
    time: "10:00",
    memory: 52,
  },
  {
    time: "10:05",
    memory: 61,
  },
  {
    time: "10:10",
    memory: 58,
  },
];

export default function MemoryChart() {
  return (
    <div
      className="
rounded-xl
border
border-zinc-800
bg-zinc-950
p-5
"
    >
      <h3 className="mb-4 text-sm text-white">Memory Usage</h3>

      <ResponsiveContainer width="100%" height={200}>
        <AreaChart data={memoryData}>
          <XAxis dataKey="time" />

          <YAxis />

          <Tooltip />

          <Area dataKey="memory" type="monotone" />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
