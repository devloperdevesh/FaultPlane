interface WorkerImpactProps {
  name: string;
  status: "healthy" | "degraded" | "failed";
}

const styles = {
  healthy: "border-emerald-500/40 bg-emerald-500/10 text-emerald-400",

  degraded: "border-orange-500/40 bg-orange-500/10 text-orange-400",

  failed: "border-red-500/40 bg-red-500/10 text-red-400",
};

export default function WorkerImpact({ name, status }: WorkerImpactProps) {
  return (
    <div
      className={`
  rounded-xl
  border
  px-4
  py-3
  transition-all
  duration-300
  hover:scale-105
  ${styles[status]}
  `}
    >
      <p
        className="
  text-sm
  font-medium
  text-white
  "
      >
        {name}
      </p>

      <p
        className="
  mt-1
  text-xs
  capitalize
  "
      >
        {status}
      </p>
    </div>
  );
}
