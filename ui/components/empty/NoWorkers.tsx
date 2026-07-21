import EmptyState from "./EmptyState";

export default function NoWorkers() {
  return (
    <EmptyState
      title="No workers detected"
      description="Runtime worker pool is currently empty"
    />
  );
}
