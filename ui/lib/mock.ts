import type {
  Worker,
  Metric,
  TelemetryEvent,
} from "./types";


export const workers: Worker[] = [

  {
    id: "worker-01",
    status: "ACTIVE",
    role: "PRIMARY",

    cpuUsage: 42,
    memoryUsage: "2GB",

    checkpointId: "v42",

    lastHeartbeat: "5s ago",
  },

];


export const metrics: Metric[] = [

  {
    name: "Latency P99",
    value: 2.3,
    unit: "ms",
  },

];


export const events: TelemetryEvent[] = [

  {
    level: "INFO",
    message: "checkpoint created",
    timestamp: "now",
    workerId: "worker-01",
  },

];