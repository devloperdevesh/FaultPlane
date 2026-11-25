import type { ReactNode } from "react";
interface RoutePageProps { eyebrow: string; title: string; description: string; children?: ReactNode; }
export default function RoutePage({ eyebrow, title, description, children }: RoutePageProps) {
 return <section className="fp-page"><header className="fp-page-heading"><div><div className="fp-eyebrow">{eyebrow}</div><h1>{title}</h1><p>{description}</p></div></header><div className="space-y-6">{children}</div></section>;
}
