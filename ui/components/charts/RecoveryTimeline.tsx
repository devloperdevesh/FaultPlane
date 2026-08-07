"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const recoveryData = [
  {
    event: "Checkpoint",
    latency: 12,
  },
  {
    event: "Failure",
    latency: 20,
  },
  {
    event: "Recovery",
    latency: 35,
  },
  {
    event: "Resume",
    latency: 18,
  },
];

export default function RecoveryTimeline() {
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
      <h3 className="mb-4 text-sm text-white">Recovery Timeline</h3>

      <ResponsiveContainer width="100%" height={240}>
        <AreaChart data={recoveryData}>
          <XAxis dataKey="event" />

          <YAxis />

          <Tooltip />

          <Area dataKey="latency" type="monotone" />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
