import DataTable from "./DataTable";
import StatusBadge from "./StatusBadge";

type WorkflowStatus =
  | "running"
  | "waiting"
  | "recovering"
  | "failed"
  | "completed";

interface Workflow {
  name: string;
  worker: string;
  status: WorkflowStatus;
  duration: string;
}

const workflows: Workflow[] = [
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
        {
          header: "Workflow",
          render: (workflow) => workflow.name,
        },

        {
          header: "Worker",
          render: (workflow) => workflow.worker,
        },

        {
          header: "State",
          render: (workflow) => <StatusBadge status={workflow.status} />,
        },

        {
          header: "Duration",
          render: (workflow) => workflow.duration,
        },
      ]}
    />
  );
}
