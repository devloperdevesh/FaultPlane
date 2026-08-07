import Panel from "./Panel";

interface Props {
  title: string;
  children: React.ReactNode;
}

export default function ChartPanel({ title, children }: Props) {
  return (
    <Panel>
      <h3
        className="
text-sm
font-semibold
text-white
mb-4
"
      >
        {title}
      </h3>

      {children}
    </Panel>
  );
}
