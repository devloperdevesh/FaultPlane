interface Props {
  id: string;
  type: "checkpoint" | "rollback";
  time: string;
}

export default function CommitNode({ id, type, time }: Props) {
  return (
    <div
      className="
  flex
  flex-col
  items-center
  "
    >
      <div
        className={`
  h-4
  w-4
  rounded-full
  border-2
  
  ${
    type === "rollback"
      ? "border-red-500 bg-red-500"
      : "border-cyan-400 bg-cyan-400"
  }
  
  shadow-lg
  `}
      ></div>

      <div
        className="
  mt-3
  text-center
  "
      >
        <p
          className="
  text-xs
  font-mono
  text-white
  "
        >
          {id}
        </p>

        <p
          className="
  text-[10px]
  text-zinc-500
  "
        >
          {time}
        </p>
      </div>
    </div>
  );
}
