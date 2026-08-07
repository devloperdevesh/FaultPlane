const capabilities = [
  {
    name: "Network",
    status: "Restricted",
  },
  {
    name: "Filesystem",
    status: "Blocked",
  },
  {
    name: "Memory",
    status: "Isolated",
  },
  {
    name: "Syscalls",
    status: "Filtered",
  },
];

export default function CapabilityMatrix() {
  return (
    <div
      className="
    rounded-xl
    border
    border-white/10
    bg-zinc-950
    p-6
    "
    >
      <h3
        className="
    text-sm
    font-semibold
    text-white
    "
      >
        Capability Matrix
      </h3>

      <div
        className="
    mt-5
    space-y-3
    "
      >
        {capabilities.map((item) => (
          <div
            key={item.name}
            className="
    flex
    justify-between
    rounded-lg
    bg-zinc-900
    p-3
    "
          >
            <span
              className="
    text-sm
    text-zinc-400
    "
            >
              {item.name}
            </span>

            <span
              className="
    text-xs
    text-emerald-400
    "
            >
              {item.status}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
