export default function PageContainer({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div
      className="
    mx-auto
    w-full
    max-w-[1600px]
    space-y-8
    "
    >
      {children}
    </div>
  );
}
