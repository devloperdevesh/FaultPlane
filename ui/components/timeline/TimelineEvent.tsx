interface Props {
  title: string;
  time: string;
}

export default function TimelineEvent({ title, time }: Props) {
  return (
    <div
      className="
    flex
    gap-4
    "
    >
      <div
        className="
    mt-1
    h-3
    w-3
    rounded-full
    bg-blue-500
    "
      />

      <div>
        <p className="text-sm text-white">{title}</p>

        <p className="text-xs text-zinc-500">{time}</p>
      </div>
    </div>
  );
}
