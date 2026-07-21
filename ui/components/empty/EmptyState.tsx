interface Props {
  title: string;
  description: string;
}

export default function EmptyState({ title, description }: Props) {
  return (
    <div
      className="
    rounded-2xl
    border
    border-white/10
    bg-zinc-950
    p-8
    text-center
    "
    >
      <h3 className="text-white font-medium">{title}</h3>

      <p className="mt-2 text-sm text-zinc-500">{description}</p>
    </div>
  );
}
