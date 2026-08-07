export default function SimulationResult() {
  return (
    <div
      className="
    mt-6
    rounded-lg
    border
    border-zinc-800
    bg-black
    p-4
    "
    >
      <h3
        className="
    text-sm
    font-semibold
    text-white
    mb-3
    "
      >
        Simulation Result
      </h3>

      <div
        className="
    space-y-2
    font-mono
    text-xs
    text-green-400
    "
      >
        <p>[INFO] Failure injected</p>

        <p>[INFO] Checkpoint located</p>

        <p>[RECOVERY] Worker restored</p>

        <p>[SUCCESS] Execution resumed</p>
      </div>
    </div>
  );
}
