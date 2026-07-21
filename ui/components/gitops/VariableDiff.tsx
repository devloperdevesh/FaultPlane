import type { VariableDiffData, VariableDiffType } from "@/lib/types";

interface Props {
  data: VariableDiffData;
}

const styles: Record<VariableDiffType, string> = {
  modified: "border-yellow-500/30 bg-yellow-500/10",

  added: "border-emerald-500/30 bg-emerald-500/10",

  removed: "border-red-500/30 bg-red-500/10",

  unchanged: "border-zinc-800 bg-zinc-900",
};

export default function VariableDiff({ data }: Props) {
  return (
    <div
      className={`
          rounded-xl
          border
          p-5
          ${styles[data.type]}
        `}
    >
      <div
        className="
          flex
          items-center
          justify-between
          "
      >
        <h3
          className="
            text-sm
            font-medium
            text-white
            "
        >
          {data.key}
        </h3>

        <span
          className="
            text-xs
            uppercase
            tracking-wide
            text-zinc-400
            "
        >
          {data.type}
        </span>
      </div>

      <div
        className="
          mt-4
          grid
          grid-cols-2
          gap-4
          font-mono
          text-sm
          "
      >
        <div>
          <p
            className="
              text-xs
              text-zinc-500
              "
          >
            Before
          </p>

          <p
            className="
              mt-2
              text-red-400
              "
          >
            - {data.before}
          </p>
        </div>

        <div>
          <p
            className="
              text-xs
              text-zinc-500
              "
          >
            After
          </p>

          <p
            className="
              mt-2
              text-emerald-400
              "
          >
            + {data.after}
          </p>
        </div>
      </div>
    </div>
  );
}
