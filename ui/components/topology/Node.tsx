interface Props {
  label: string;
  type: string;
}

export default function Node({ label, type }: Props) {
  return (
    <div
      className="
    rounded-xl
    border
    border-white/10
    bg-zinc-900
    px-5
    py-4
    text-center
    "
    >
      <p className="text-xs text-zinc-500">{type}</p>

      <p className="mt-2 text-sm text-white">{label}</p>
    </div>
  );
}
