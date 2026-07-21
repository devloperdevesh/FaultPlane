import WorkerTable from "@/components/workers/WorkerTable";

export default function WorkersPage() {
  return (
    <div className="space-y-6">
      <h1
        className="
text-3xl
font-semibold
text-white
"
      >
        Workers
      </h1>

      <p className="text-zinc-500">12 Active Runtime Nodes</p>

      <WorkerTable />
    </div>
  );
}
