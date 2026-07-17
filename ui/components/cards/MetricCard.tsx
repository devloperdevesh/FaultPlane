interface MetricCardProps {
  title: string;
  value: string;
  description?: string;
}

export default function MetricCard({
  title,
  value,
  description,
}: MetricCardProps) {
  return (
    <div
      className="
        rounded-xl
        border
        border-zinc-800
        bg-zinc-950
        p-5
        "
    >
      <p className="text-xs uppercase text-zinc-500">{title}</p>

      <h3 className="mt-2 text-2xl font-semibold text-white">{value}</h3>

      {description && (
        <p className="mt-2 text-xs text-zinc-500">{description}</p>
      )}
    </div>
  );
}
