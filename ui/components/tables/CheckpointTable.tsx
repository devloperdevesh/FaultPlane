"use client";

import DataTable from "./DataTable";

type Checkpoint = {
  id: string;
  worker: string;
  size: string;
  storage: string;
  created: string;
};

export default function CheckpointTable() {
  const checkpoints: Checkpoint[] = [];

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
      ]}    />
  );
}