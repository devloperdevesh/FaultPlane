import DataTable from "./DataTable";
import StatusBadge from "./StatusBadge";

type WorkerStatus = "healthy" | "recovering" | "offline";

interface Worker {
  id: string;
  role: string;
  status: WorkerStatus;
  cpu: string;
  memory: string;
  requests: number;
}

const workers: Worker[] = [
  {
    id: "worker-01",
    role: "Primary",
    status: "healthy",
    cpu: "28%",
    memory: "1.7 GB",
    requests: 18421,
  },
  {
    id: "worker-02",
    role: "Standby",
    status: "recovering",
    cpu: "17%",
    memory: "1.1 GB",
    requests: 8921,
  },
];

export default function WorkersTable() {
  return (
    <DataTable
      title="Workers"
      data={workers}
      columns={[
        {
          header: "Worker",
          render: (worker) => worker.id,
        },

        {
          header: "Role",
          render: (worker) => worker.role,
        },

        {
          header: "Status",
          render: (worker) => <StatusBadge status={worker.status} />,
        },

        {
          header: "CPU",
          render: (worker) => worker.cpu,
        },

        {
          header: "Memory",
          render: (worker) => worker.memory,
        },

        {
          header: "Requests",
          render: (worker) => worker.requests,
        },
      ]}
    />
  );
}
