import RoutePage from "@/components/layout/RoutePage";
import SpeculativePath from "@/components/workflows/SpeculativePath";
import ExecutionTree from "@/components/experimental/speculative/ExecutionTree";

export default function SpeculativePage() {
  return (
    <RoutePage
      eyebrow="Experiment"
      title="Speculative Paths"
      description="Explore alternative execution branches and speculative recovery behavior in a dedicated lab view."
    >
      <div className="grid gap-6 xl:grid-cols-2">
        <SpeculativePath />
        <ExecutionTree />
      </div>
    </RoutePage>
  );
}
