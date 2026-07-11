import "./globals.css";

import AppShell from "@/components/layout/AppShell";

export const metadata = {
  title: "FaultPlane",
  description: "FaultPlane Operations Dashboard",
};

function SidebarPlaceholder() {
  return <div className="p-6">Sidebar</div>;
}

function NavbarPlaceholder() {
  return <div className="flex h-full items-center px-6">FaultPlane</div>;
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <AppShell
          sidebar={<SidebarPlaceholder />}
          navbar={<NavbarPlaceholder />}
        >
          {children}
        </AppShell>
      </body>
    </html>
  );
}
