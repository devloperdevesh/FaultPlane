interface WorkerNodeProps {
  id: string;
  status: string;
}

export default function WorkerNode({ id, status }: WorkerNodeProps) {
  return (
    <div
      className="
        rounded-xl
        border
        border-purple-500/30
        bg-zinc-950
        px-6
        py-4
        text-center
        "
    >
      <p className="text-sm font-semibold text-white">{id}</p>

      <span
        className="
          mt-2
          inline-block
          rounded-full
          bg-purple-500/10
          px-3
          py-1
          text-xs
          text-purple-400
          "
      >
        {status}
      </span>
    </div>
  );
}
