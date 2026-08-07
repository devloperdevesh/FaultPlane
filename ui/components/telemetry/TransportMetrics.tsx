export default function TransportMetrics() {
  const data = [
    ["Connections", "1.2K"],
    ["Packets/s", "84K"],
    ["Latency", "4.8ms"],
    ["Dropped", "0.02%"],
  ];

  return (
    <div
      className="
    grid
    gap-4
    md:grid-cols-4
    "
    >
      {data.map((item) => (
        <div
          key={item[0]}
          className="
    rounded-xl
    border
    border-white/10
    bg-zinc-950
    p-4
    "
        >
          <p
            className="
    text-xs
    text-zinc-500
    "
          >
            {item[0]}
          </p>

          <p
            className="
    mt-2
    font-mono
    text-white
    "
          >
            {item[1]}
          </p>
        </div>
      ))}
    </div>
  );
}
