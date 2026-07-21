interface TimelineNodeProps {
  title: string;
  timestamp: string;
  status?: "active" | "completed" | "failed";
}

const statusStyles = {
  active: "bg-blue-500 border-blue-400",
  completed: "bg-emerald-500 border-emerald-400",
  failed: "bg-red-500 border-red-400",
};

export default function TimelineNode({
  title,
  timestamp,
  status = "completed",
}: TimelineNodeProps) {
  return (
    <div className="flex flex-col items-center">
      <div
        className={`
          h-4
          w-4
          rounded-full
          border-2
          ${statusStyles[status]}
          `}
      />

      <p
        className="
          mt-3
          text-sm
          font-medium
          text-white
          "
      >
        {title}
      </p>

      <span
        className="
          mt-1
          text-xs
          text-zinc-500
          "
      >
        {timestamp}
      </span>
    </div>
  );
}
