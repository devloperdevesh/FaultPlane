import TimelineNode from "./TimelineNode";

const checkpoints = [
  {
    title: "checkpoint-01",
    timestamp: "10:01:12",
  },
  {
    title: "checkpoint-02",
    timestamp: "10:05:24",
  },
  {
    title: "checkpoint-03",
    timestamp: "10:10:42",
  },
  {
    title: "checkpoint-04",
    timestamp: "10:15:18",
  },
];

export default function CheckpointTimeline() {
  return (
    <div
      className="
      rounded-xl
      border
      border-zinc-800
      bg-zinc-950
      p-6
      "
    >
      <div className="mb-8">
        <h2 className="text-sm font-semibold text-white">Checkpoint History</h2>

        <p className="text-xs text-zinc-500">Runtime state snapshots</p>
      </div>

      <div
        className="
        flex
        items-start
        justify-between
        "
      >
        {checkpoints.map((checkpoint, index) => (
          <div key={checkpoint.title} className="flex flex-1 items-center">
            <TimelineNode
              title={checkpoint.title}
              timestamp={checkpoint.timestamp}
            />

            {index !== checkpoints.length - 1 && (
              <div
                className="
                mx-4
                h-px
                flex-1
                bg-zinc-700
                "
              />
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
