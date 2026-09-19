import DataTable from "./DataTable";

type Workflow = {
  name: string;
  worker: string;
  status: string;
  duration: string;
};

export default function WorkflowTable() {
  const workflows: Workflow[] = [];

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
          render: (workflow) => workflow.status,
        },
        {
          header: "Duration",
          render: (workflow) => workflow.duration,
        },
      ]}
    />
  );
}
