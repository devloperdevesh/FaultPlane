const events = [
  {
    time: "10:32:01",
    event: "Worker failure detected",
  },
  {
    time: "10:32:02",
    event: "Blast radius contained",
  },
  {
    time: "10:32:04",
    event: "Checkpoint restored",
  },
  {
    time: "10:32:06",
    event: "Recovery completed",
  },
];

export default function ImpactTimeline() {
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
      <h3
        className="
    text-sm
    font-semibold
    text-white
    "
      >
        Impact Timeline
      </h3>

      <div
        className="
    mt-6
    space-y-5
    "
      >
        {events.map((item) => (
          <div
            key={item.time}
            className="
    flex
    gap-4
    "
          >
            <div
              className="
    font-mono
    text-xs
    text-zinc-500
    "
            >
              {item.time}
            </div>

            <div
              className="
    text-sm
    text-white
    "
            >
              {item.event}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
