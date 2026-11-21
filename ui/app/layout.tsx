import type { ReactNode } from "react";
import "./globals.css";
import AppFrame from "@/components/shell/AppFrame";
export const metadata = { title: "FaultPlane — AI Runtime Control Plane", description: "Runtime resilience for long-running AI workloads. Detect infrastructure signals, evaluate policies and verify recovery." };
export default function RootLayout({ children }: { children: ReactNode }) {
 return <html lang="en" suppressHydrationWarning><body><AppFrame>{children}</AppFrame></body></html>;
}
