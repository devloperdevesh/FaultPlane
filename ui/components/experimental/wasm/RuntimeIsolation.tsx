const permissions = ["Filesystem", "Network", "Memory", "Permissions"];

export default function RuntimeIsolation() {
  return (
    <div
      className="
    grid
    grid-cols-2
    gap-3
    "
    >
      {permissions.map((item) => (
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
          {item}
        </div>
      ))}
    </div>
  );
}
