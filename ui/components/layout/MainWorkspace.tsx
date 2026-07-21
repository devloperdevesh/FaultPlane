interface Props {
  children: React.ReactNode;
}

export default function MainWorkspace({ children }: Props) {
  return (
    <main
      className="
   flex-1
   overflow-y-auto
   bg-black
   p-8
   "
    >
      {children}
    </main>
  );
}
