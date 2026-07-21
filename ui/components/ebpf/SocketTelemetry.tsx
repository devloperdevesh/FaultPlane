import Panel from "@/components/ui/Panel";

const data = [
  {
    name: "TCP Connections",
    value: "248",
  },
  {
    name: "QUIC Streams",
    value: "64",
  },
  {
    name: "Packets/sec",
    value: "84K",
  },
  {
    name: "Socket Migration",
    value: "12",
  },
];

export default function SocketTelemetry() {
  return (
    <Panel>
      <h2 className="text-white mb-4 font-semibold">Socket Telemetry</h2>

      <div className="grid grid-cols-2 gap-4">
        {data.map((item) => (
          <div
            key={item.name}
            className="
rounded-lg
bg-zinc-900
p-4
"
          >
            <p className="text-xs text-zinc-500">{item.name}</p>

            <p className="text-xl text-white mt-2">{item.value}</p>
          </div>
        ))}
      </div>
    </Panel>
  );
}
