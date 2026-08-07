import TimelineEvent from "./TimelineEvent";

const events = [
  ["Checkpoint Created", "10:01:02"],

  ["Worker Failed", "10:01:14"],

  ["Recovery Started", "10:01:15"],

  ["Execution Resumed", "10:01:28"],
];

export default function RecoveryTimeline() {
  return (
    <div
      className="
space-y-6
rounded-2xl
border
border-white/10
bg-zinc-950
p-6
"
    >
      <h2
        className="
text-sm
font-semibold
text-white
"
      >
        Recovery Timeline
      </h2>

      <div
        className="
space-y-5
"
      >
        {events.map(([title, time]) => (
          <TimelineEvent key={title} title={title} time={time} />
        ))}
      </div>
    </div>
  );
}
