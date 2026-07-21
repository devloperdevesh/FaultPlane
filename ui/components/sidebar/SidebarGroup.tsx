import type { ReactNode } from "react";

interface Props {
  title: string;
  children: ReactNode;
}

export default function SidebarGroup({ title, children }: Props) {
  return (
    <section className="mb-6">
      <h3
        className="
mb-2
px-3
text-xs
uppercase
tracking-widest
text-zinc-500
"
      >
        {title}
      </h3>

      <div className="space-y-1">{children}</div>
    </section>
  );
}
