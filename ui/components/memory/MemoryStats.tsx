interface MemoryStatsProps {
  title: string;
  value: string;
  status?: "default" | "success" | "warning";
}

export default function MemoryStats({
  title,
  value,
  status = "default",
}: MemoryStatsProps) {
  const valueClass = {
    default: "text-white",
    success: "text-emerald-400",
    warning: "text-amber-400",
  };

  return (
    <div
      className="
          rounded-xl
          border
          border-zinc-800
          bg-zinc-900/60
          p-4
          transition-all
          duration-200
          hover:border-zinc-700
          hover:bg-zinc-900
        "
    >
      <p
        className="
            text-xs
            font-medium
            uppercase
            tracking-wide
            text-zinc-500
          "
      >
        {title}
      </p>

      <p
        className={`
            mt-2
            font-mono
            text-lg
            font-semibold
            ${valueClass[status]}
          `}
      >
        {value}
      </p>
    </div>
  );
}
