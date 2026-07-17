import DataTable from "./DataTable";
import StatusBadge from "./StatusBadge";

const workflows = [
  {
    name: "Chat Agent",
    worker: "worker-01",
    status: "running",
    duration: "18m",
  },
  {
    name: "RAG Index",
    worker: "worker-02",
    status: "waiting",
    duration: "3m",
  },
];

export default function WorkflowTable() {
  return (
    <DataTable
      title="Workflows"
      data={workflows}
      columns={[
        { header: "Workflow", render: (w) => w.name },
        { header: "Worker", render: (w) => w.worker },
        {
          header: "State",
          render: (w) => <StatusBadge status={w.status} />,
        },
        { header: "Duration", render: (w) => w.duration },
      ]}
    />
  );
}
