import TimelineNode from "./TimelineNode";

const events = [
  {
    title: "Runtime Started",
    timestamp: "09:00",
  },
  {
    title: "Checkpoint Created",
    timestamp: "09:05",
  },
  {
    title: "Worker Recovery",
    timestamp: "09:12",
  },
];

export default function GitTimeline() {
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
      <h2
        className="
        mb-6
        text-sm
        font-semibold
        text-white
        "
      >
        Execution Timeline
      </h2>

      <div
        className="
        flex
        items-center
        justify-between
        "
      >
        {events.map((event) => (
          <TimelineNode
            key={event.title}
            title={event.title}
            timestamp={event.timestamp}
          />
        ))}
      </div>
    </div>
  );
}
