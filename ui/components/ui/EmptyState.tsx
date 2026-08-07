interface Props {
  title: string;
  description?: string;
}

export default function EmptyState({ title, description }: Props) {
  return (
    <div
      className="
    rounded-xl
    border
    border-dashed
    border-zinc-800
    p-8
    text-center
    "
    >
      <h3 className="text-white">{title}</h3>

      <p className="mt-2 text-sm text-zinc-500">{description}</p>
    </div>
  );
}
