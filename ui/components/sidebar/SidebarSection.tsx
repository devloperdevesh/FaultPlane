import { ReactNode } from "react";

interface SidebarSectionProps {
  title: string;
  children: ReactNode;
}

export default function SidebarSection({
  title,
  children,
}: SidebarSectionProps) {
  return (
    <section className="mb-8">
      <h2
        className="
        mb-3
        px-4
        text-[11px]
        font-semibold
        uppercase
        tracking-[0.22em]
        text-zinc-500
        "
      >
        {title}
      </h2>

      <div className="space-y-1">{children}</div>
    </section>
  );
}
