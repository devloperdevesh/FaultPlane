import DataTable from "./DataTable";

const checkpoints = [
  {
    id: "cp-1023",
    worker: "worker-01",
    size: "18 MB",
    created: "2 min ago",
    storage: "Memory",
  },
  {
    id: "cp-1024",
    worker: "worker-02",
    size: "25 MB",
    created: "6 min ago",
    storage: "Disk",
  },
];

export default function CheckpointTable() {
  return (
    <DataTable
      title="Checkpoints"
      data={checkpoints}
      columns={[
        { header: "ID", render: (c) => c.id },
        { header: "Worker", render: (c) => c.worker },
        { header: "Size", render: (c) => c.size },
        { header: "Storage", render: (c) => c.storage },
        { header: "Created", render: (c) => c.created },
      ]}
    />
  );
}
