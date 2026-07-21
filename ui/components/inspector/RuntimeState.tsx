"use client";

interface Props {
  state: "RUNNING" | "FAILED" | "RECOVERING";
}

export default function RuntimeState({ state }: Props) {
  const colors = {
    RUNNING: "text-emerald-400 bg-emerald-500/10",

    FAILED: "text-red-400 bg-red-500/10",

    RECOVERING: "text-yellow-400 bg-yellow-500/10",
  };

  return (
    <span
      className={`
rounded-full
px-3
py-1
text-xs
font-medium
${colors[state]}
`}
    >
      {state}
    </span>
  );
}
