interface ConnectionProps {
  label?: string;
}

export default function Connection({ label }: ConnectionProps) {
  return (
    <div
      className="
        flex
        flex-col
        items-center
        "
    >
      <div
        className="
          h-12
          w-px
          bg-zinc-700
          "
      />

      {label && (
        <span
          className="
            text-xs
            text-zinc-500
            "
        >
          {label}
        </span>
      )}
    </div>
  );
}
