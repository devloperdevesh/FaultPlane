interface Props {
  state: "stable" | "migrating" | "restored" | "empty";
}

const styles = {
  stable: "bg-cyan-400/80 shadow-cyan-400/40",

  migrating: "bg-yellow-400 animate-pulse shadow-yellow-400/60",

  restored: "bg-emerald-400 shadow-emerald-400/40",

  empty: "bg-zinc-800",
};

export default function MemoryPage({ state }: Props) {
  return (
    <div
      className={`
        h-8
        w-8
        rounded-md
        shadow-lg
        transition-all
        duration-500
        ${styles[state]}
        `}
    />
  );
}
