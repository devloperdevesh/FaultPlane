"use client";

interface Props {
  cpu: string;
  memory: string;
}

export default function WorkerHealth({ cpu, memory }: Props) {
  return (
    <div
      className="
space-y-1
font-mono
text-xs
"
    >
      <p className="text-zinc-400">
        CPU
        <span className="ml-2 text-white">{cpu}</span>
      </p>

      <p className="text-zinc-400">
        Memory
        <span className="ml-2 text-white">{memory}</span>
      </p>
    </div>
  );
}
