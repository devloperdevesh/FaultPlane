import { ReactNode } from "react";

interface PanelProps {
  children: ReactNode;
  className?: string;
}

export default function Panel({ children, className = "" }: PanelProps) {
  return (
    <section
      className={`
        rounded-xl
        border
        border-zinc-800
        bg-zinc-950
        p-6
        transition-all
        duration-200
        hover:border-zinc-700
        ${className}
      `}
    >
      {children}
    </section>
  );
}
