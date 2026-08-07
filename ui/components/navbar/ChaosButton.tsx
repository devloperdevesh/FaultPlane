"use client";

import { Flame } from "lucide-react";
import { useState } from "react";

export default function ChaosButton() {
  const [loading, setLoading] = useState(false);

  function handleChaos() {
    setLoading(true);

    setTimeout(() => {
      setLoading(false);
    }, 1500);
  }

  return (
    <button
      onClick={handleChaos}
      disabled={loading}
      className="
        flex items-center gap-2
        rounded-lg
        bg-gradient-to-r
        from-red-600
        to-orange-500
        px-4
        py-2
        text-sm
        font-semibold
        text-white
        shadow-lg
        transition
        hover:scale-105
        disabled:opacity-50
      "
    >
      <Flame className="h-4 w-4" />

      {loading ? "Injecting..." : "Simulate Failure"}
    </button>
  );
}
