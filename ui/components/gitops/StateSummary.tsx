interface SummaryStat {
  title: string;
  value: string;
}

const stats: SummaryStat[] = [
  {
    title: "Variables Changed",
    value: "2",
  },
  {
    title: "Memory Restored",
    value: "100%",
  },
  {
    title: "Checkpoint",
    value: "cp-1024",
  },
];

export default function StateSummary() {
  return (
    <div className="grid gap-4 md:grid-cols-3">
      {stats.map((stat) => (
        <div
          key={stat.title}
          className="
              rounded-xl
              border
              border-white/10
              bg-zinc-900
              p-5
              transition-colors
              hover:border-emerald-500/30
            "
        >
          <p className="text-xs uppercase tracking-wide text-zinc-500">
            {stat.title}
          </p>

          <p className="mt-2 font-mono text-lg font-semibold text-white">
            {stat.value}
          </p>
        </div>
      ))}
    </div>
  );
}
