"use client";

const steps = [
  "Snapshot Selected",
  "State Restored",
  "Worker Restarted",
  "Traffic Shifted",
];

export default function RollbackDiff() {
  return (
    <div
      className="
rounded-2xl
border
border-white/10
bg-zinc-950
p-6
"
    >
      <h2
        className="
text-sm
font-semibold
text-white
"
      >
        Rollback Execution
      </h2>

      <div
        className="
mt-5
space-y-3
"
      >
        {steps.map((step, index) => (
          <div
            key={step}
            className="
flex
items-center
gap-4
"
          >
            <div
              className="
h-8
w-8
rounded-full
bg-emerald-500/20
text-emerald-400
flex
items-center
justify-center
text-xs
"
            >
              {index + 1}
            </div>

            <p
              className="
text-sm
text-zinc-300
"
            >
              {step}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
