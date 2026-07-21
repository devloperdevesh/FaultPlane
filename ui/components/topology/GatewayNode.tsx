interface GatewayNodeProps {
  name: string;
  status: string;
}

export default function GatewayNode({ name, status }: GatewayNodeProps) {
  return (
    <div
      className="
        rounded-xl
        border
        border-blue-500/30
        bg-zinc-950
        px-6
        py-4
        text-center
        "
    >
      <p className="text-sm font-semibold text-white">{name}</p>

      <span
        className="
          mt-2
          inline-block
          rounded-full
          bg-emerald-500/10
          px-3
          py-1
          text-xs
          text-emerald-400
          "
      >
        {status}
      </span>
    </div>
  );
}
