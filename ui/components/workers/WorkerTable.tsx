"use client";

import WorkerStatus from "./WorkerStatus";
import WorkerHealth from "./WorkerHealth";

const workers = [
  {
    id: "worker-01",
    status: "ACTIVE",
    cpu: "42%",
    memory: "2GB",
    checkpoint: "v42",
    role: "PRIMARY",
  },

  {
    id: "worker-02",
    status: "ACTIVE",
    cpu: "31%",
    memory: "1.4GB",
    checkpoint: "v41",
    role: "STANDBY",
  },

  {
    id: "worker-03",
    status: "RECOVERING",
    cpu: "12%",
    memory: "900MB",
    checkpoint: "v39",
    role: "STANDBY",
  },
] as const;

export default function WorkerTable() {
  return (
    <div
      className="
overflow-hidden
rounded-2xl
border
border-white/10
bg-zinc-950
"
    >
      <table
        className="
w-full
text-sm
"
      >
        <thead
          className="
border-b
border-white/10
text-xs
uppercase
text-zinc-500
"
        >
          <tr>
            <th className="px-6 py-4 text-left">Worker</th>

            <th>Status</th>

            <th>Health</th>

            <th>Checkpoint</th>

            <th>Role</th>
          </tr>
        </thead>

        <tbody>
          {workers.map((worker) => (
            <tr
              key={worker.id}
              className="
border-b
border-white/5
hover:bg-white/5
transition
"
            >
              <td
                className="
px-6
py-4
font-mono
text-white
"
              >
                {worker.id}
              </td>

              <td>
                <WorkerStatus status={worker.status} />
              </td>

              <td>
                <WorkerHealth cpu={worker.cpu} memory={worker.memory} />
              </td>

              <td
                className="
font-mono
text-zinc-300
"
              >
                {worker.checkpoint}
              </td>

              <td
                className="
text-zinc-400
"
              >
                {worker.role}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
