"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const data = [
  {
    time: "10:00",
    p50: 80,
    p95: 160,
    p99: 240,
  },
  {
    time: "10:05",
    p50: 90,
    p95: 180,
    p99: 260,
  },
  {
    time: "10:10",
    p50: 70,
    p95: 150,
    p99: 220,
  },
  {
    time: "10:15",
    p50: 110,
    p95: 210,
    p99: 320,
  },
];

export default function LatencyChart() {
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
      <h3 className="mb-4 text-sm font-medium text-white">
        Latency Distribution
      </h3>

      <ResponsiveContainer width="100%" height={260}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />

          <XAxis dataKey="time" />

          <YAxis />

          <Tooltip />

          <Line type="monotone" dataKey="p50" strokeWidth={2} />

          <Line type="monotone" dataKey="p95" strokeWidth={2} />

          <Line type="monotone" dataKey="p99" strokeWidth={2} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
