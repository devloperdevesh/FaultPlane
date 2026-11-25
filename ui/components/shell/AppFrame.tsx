"use client";
import type { ReactNode } from "react";
import { usePathname } from "next/navigation";
import Sidebar from "@/components/sidebar/Sidebar";
import Topbar from "./Topbar";
export default function AppFrame({ children }: { children: ReactNode }) {
 const pathname=usePathname();
 if (["/login","/signup","/forgot-password","/checkout"].includes(pathname)) return <>{children}</>;
 return <div className="fp-app"><Sidebar/><div className="fp-main"><Topbar/><main className="fp-content">{children}</main><footer className="fp-footer"><span><i className="fp-dot"/> FaultPlane control plane</span><span>Runtime connection status is reported by the connected agent</span><span>© 2026 FaultPlane</span></footer></div></div>;
}
