"use client";

interface Props {
  title: string;
  probability: string;
  selected?: boolean;
}

export default function PredictionNode({
  title,
  probability,
  selected,
}: Props) {
  return (
    <div
      className={`
rounded-xl
border
p-4
transition
${
  selected
    ? "border-emerald-500/40 bg-emerald-500/10"
    : "border-white/10 bg-zinc-900"
}
`}
    >
      <p className="text-sm text-white">{title}</p>

      <p className="mt-2 text-xs text-zinc-400">{probability}% success</p>
    </div>
  );
}
