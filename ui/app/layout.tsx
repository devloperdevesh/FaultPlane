import type { ReactNode } from "react";

import "./globals.css";

import Sidebar from "@/components/sidebar/Sidebar";

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body className="h-screen overflow-hidden bg-[#060914] text-slate-200 antialiased">
        <div
          className="flex h-screen overflow-hidden"
          style={{
            backgroundImage:
              "radial-gradient(circle at top, rgba(34, 197, 94, 0.10), transparent 30%), linear-gradient(180deg, #060914 0%, #05070f 100%)",
          }}
        >
          <Sidebar />

          <main className="min-w-0 flex-1 overflow-y-auto">
            <div className="min-h-screen px-4 py-6 sm:px-6 lg:px-8">
              {children}
            </div>
          </main>
        </div>
      </body>
    </html>
  );
}
