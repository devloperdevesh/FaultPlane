import "./globals.css";

export const metadata = {
  title: "FaultPlane",
  description: "FaultPlane Operations Dashboard",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
