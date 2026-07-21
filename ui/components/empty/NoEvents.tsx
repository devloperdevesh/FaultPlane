import EmptyState from "./EmptyState";

export default function NoEvents() {
  return (
    <EmptyState
      title="No recovery events detected"
      description="Your system is operating normally"
    />
  );
}
