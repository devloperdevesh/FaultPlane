import Panel from "@/components/ui/Panel";

const events = [
  "[INFO] xdp program attached",

  "[INFO] socket state captured",

  "[RECOVERY] connection migrated",

  "[CHECKPOINT] memory snapshot created",
];

export default function KernelEvents() {
  return (
    <Panel>
      <h2 className="text-white mb-4">Kernel Events</h2>

      <div
        className="
font-mono
text-xs
space-y-3
text-green-400
"
      >
        {events.map((event) => (
          <div key={event}>{event}</div>
        ))}
      </div>
    </Panel>
  );
}
