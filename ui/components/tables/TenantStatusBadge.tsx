interface Props {
  status: "healthy" | "warning" | "critical";
}

const config = {
  healthy: {
    label: "Healthy",
    style: "text-emerald-400 bg-emerald-500/10 border-emerald-500/30",
  },

  warning: {
    label: "Warning",
    style: "text-orange-400 bg-orange-500/10 border-orange-500/30",
  },

  critical: {
    label: "Critical",
    style: "text-red-400 bg-red-500/10 border-red-500/30",
  },
};

export default function TenantStatusBadge({ status }: Props) {
  const item = config[status];

  return (
    <span
      className={`
  rounded-full
  border
  px-3
  py-1
  text-xs
  font-medium
  ${item.style}
  `}
    >
      {item.label}
    </span>
  );
}
