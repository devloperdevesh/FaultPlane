interface Props {
  retry?: () => void;
}

export default function ErrorState({ retry }: Props) {
  return (
    <div
      className="
    rounded-2xl
    border
    border-red-500/30
    bg-red-500/10
    p-8
    "
    >
      <h3 className="text-red-400 font-semibold">
        Runtime telemetry unavailable
      </h3>

      <p className="mt-2 text-sm text-zinc-400">
        Unable to connect with backend services.
      </p>

      <button
        onClick={retry}
        className="
    mt-4
    rounded-lg
    bg-red-500
    px-4
    py-2
    text-sm
    text-white
    "
      >
        Retry connection
      </button>
    </div>
  );
}
