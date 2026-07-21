import AgentInspector from "@/components/inspector/AgentInspector";

export default function InspectorPanel() {
  return (
    <aside
      className="
w-80
border-l
border-zinc-800
bg-zinc-950
p-5
"
    >
      <AgentInspector />
    </aside>
  );
}
