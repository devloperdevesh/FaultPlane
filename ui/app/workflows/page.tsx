import RoutePage from "@/components/layout/RoutePage";
import WorkflowTable from "@/components/tables/WorkflowTable";
import SpeculativePath from "@/components/workflows/SpeculativePath";
import ExecutionTree from "@/components/experimental/speculative/ExecutionTree";

export default function WorkflowsPage() {
  return (
    <RoutePage
      eyebrow="Orchestration"
      title="Workflows"
      description="Track workflow execution, recovery branches, and speculative paths across the runtime control plane."
    >
      <div className="space-y-6">
        <WorkflowTable />

        <div className="grid gap-6 xl:grid-cols-2">
          <SpeculativePath />
          <ExecutionTree />
        </div>
      </div>
    </RoutePage>
  );
}
