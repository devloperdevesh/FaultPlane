import type { ReactNode } from "react";

import PageContainer from "./PageContainer";

interface RoutePageProps {
  eyebrow: string;
  title: string;
  description: string;
  children?: ReactNode;
}

export default function RoutePage({
  eyebrow,
  title,
  description,
  children,
}: RoutePageProps) {
  return (
    <PageContainer>
      <section className="space-y-6">
        <header className="space-y-3 border-b border-white/10 pb-6">
          <p className="text-xs font-semibold uppercase tracking-[0.28em] text-emerald-400/90">
            {eyebrow}
          </p>

          <h1 className="text-3xl font-semibold tracking-tight text-white">
            {title}
          </h1>

          <p className="max-w-3xl text-sm leading-6 text-zinc-400">
            {description}
          </p>
        </header>

        {children}
      </section>
    </PageContainer>
  );
}
