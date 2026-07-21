const kernelStatus = [
  {
    name: "eBPF Hook",
    value: "ACTIVE",
  },
  {
    name: "Socket Watch",
    value: "ACTIVE",
  },
  {
    name: "Failover Engine",
    value: "READY",
  },
];

export default function KernelMonitor() {
  return (
    <div
      className="
        rounded-xl
        border
        border-zinc-800
        bg-black
        p-5
        "
    >
      <h2
        className="
          mb-5
          text-sm
          font-semibold
          text-white
          "
      >
        Kernel Monitor
      </h2>

      <div className="space-y-3">
        {kernelStatus.map((item) => (
          <div
            key={item.name}
            className="
              flex
              justify-between
              rounded-lg
              bg-zinc-900
              px-3
              py-2
              "
          >
            <span
              className="
                font-mono
                text-xs
                text-zinc-400
                "
            >
              {item.name}
            </span>

            <span
              className="
                font-mono
                text-xs
                text-emerald-400
                "
            >
              {item.value}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
