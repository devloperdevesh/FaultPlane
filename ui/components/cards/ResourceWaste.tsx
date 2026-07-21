export default function ResourceWaste() {
  const items = [
    "Idle CPU detected",
    "Unused memory reservation",
    "Over-provisioned workers",
  ];

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
        Waste Detection
      </h3>

      <div
        className="
    mt-4
    space-y-3
    "
      >
        {items.map((item) => (
          <div
            key={item}
            className="
    rounded-lg
    bg-zinc-900
    p-3
    text-sm
    text-zinc-300
    "
          >
            ⚠ {item}
          </div>
        ))}
      </div>
    </div>
  );
}
