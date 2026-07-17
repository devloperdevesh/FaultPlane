import DataTable from "./DataTable";
import StatusBadge from "./StatusBadge";

const workers = [
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
        { header: "Worker", render: (w) => w.id },
        { header: "Role", render: (w) => w.role },
        {
          header: "Status",
          render: (w) => <StatusBadge status={w.status} />,
        },
        { header: "CPU", render: (w) => w.cpu },
        { header: "Memory", render: (w) => w.memory },
        { header: "Requests", render: (w) => w.requests },
      ]}
    />
  );
}
