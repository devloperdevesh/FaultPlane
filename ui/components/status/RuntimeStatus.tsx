interface StatusItem {
  name: string;
  status: string;
}

const statusItems: StatusItem[] = [
  {
    name: "Gateway",
    status: "Healthy",
  },
  {
    name: "Storage",
    status: "Connected",
  },
  {
    name: "Workers",
    status: "2 Active",
  },
];

export default function RuntimeStatus() {
  return (
    <div
      className="
        rounded-xl
        border
        border-zinc-800
        bg-zinc-950
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
        Runtime Status
      </h2>

      <div className="space-y-3">
        {statusItems.map((item) => (
          <div
            key={item.name}
            className="
              flex
              justify-between
              "
          >
            <span className="text-sm text-zinc-400">{item.name}</span>

            <span
              className="
                text-sm
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
