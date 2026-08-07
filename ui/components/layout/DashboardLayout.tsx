import InspectorPanel from "./InspectorPanel";

interface Props {
  children: React.ReactNode;

  inspector?: React.ReactNode;
}

export default function DashboardLayout({ children, inspector }: Props) {
  return (
    <div
      className="
flex
h-full
"
    >
      <section
        className="
flex-1
"
      >
        {children}
      </section>

      <InspectorPanel>{inspector}</InspectorPanel>
    </div>
  );
}
