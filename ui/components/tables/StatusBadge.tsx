interface StatusBadgeProps {
  status:
    | "healthy"
    | "running"
    | "recovering"
    | "waiting"
    | "failed"
    | "offline"
    | "completed";
}

const styles = {
  healthy: "bg-emerald-500/10 text-emerald-400",
  running: "bg-blue-500/10 text-blue-400",
  recovering: "bg-amber-500/10 text-amber-400",
  waiting: "bg-zinc-700 text-zinc-300",
  failed: "bg-red-500/10 text-red-400",
  offline: "bg-red-500/10 text-red-400",
  completed: "bg-emerald-500/10 text-emerald-400",
};

export default function StatusBadge({ status }: StatusBadgeProps) {
  return (
    <span
      className={`rounded-full px-2 py-1 text-xs font-medium capitalize ${styles[status]}`}
    >
      {status}
    </span>
  );
}
