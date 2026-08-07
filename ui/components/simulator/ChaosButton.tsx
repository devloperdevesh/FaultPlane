"use client";

import { useState } from "react";
import ChaosModal from "./ChaosModal";

export default function ChaosButton() {
  const [open, setOpen] = useState(false);

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="
        rounded-lg
        bg-red-600
        px-4
        py-2
        text-sm
        font-semibold
        text-white
        transition
        hover:bg-red-700
        "
      >
        Inject Failure
      </button>

      {open && <ChaosModal onClose={() => setOpen(false)} />}
    </>
  );
}
