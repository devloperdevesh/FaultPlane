"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  ResponsiveContainer,
  Tooltip,
} from "recharts";

const cpuData = [
  {
    time: "10:00",
    cpu: 35,
  },
  {
    time: "10:05",
    cpu: 48,
  },
  {
    time: "10:10",
    cpu: 42,
  },
];

export default function CpuChart() {
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
      <h3 className="text-sm text-white mb-4">CPU Usage</h3>

      <ResponsiveContainer width="100%" height={200}>
        <LineChart data={cpuData}>
          <XAxis dataKey="time" />

          <YAxis />

          <Tooltip />

          <Line dataKey="cpu" type="monotone" strokeWidth={2} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
