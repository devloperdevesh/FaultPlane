"use client";

import { FaGithub } from "react-icons/fa";

export default function GithubShortcut() {
  return (
    <a
      href="https://github.com/devloperdevesh/FaultPlane"
      target="_blank"
      rel="noopener noreferrer"
      aria-label="Open GitHub repository"
      className="
        flex
        h-9
        w-9
        items-center
        justify-center
        rounded-xl
        border
        border-white/10
        bg-zinc-950/70
        text-zinc-400
        transition-all
        duration-200
        hover:border-white/20
        hover:bg-white/10
        hover:text-white
      "
    >
      <FaGithub size={18} />
    </a>
  );
}
