export interface DashboardMetrics {
  requests: number;
  workers: number;
  recoveries: number;
  checkpoints: number;
  latency: number;
  cpu: number;
  memory: number;
  updatedAt: string;
}

export interface Metric {
  name: string;
  value: number;
  unit?: string;
  timestamp?: string;
}

export type WorkerStatus =
  | "ACTIVE"
  | "FAILED"
  | "RECOVERING"
  | "STOPPED";

export type WorkerRole =
  | "PRIMARY"
  | "STANDBY";

export interface Worker {
  id: string;
  status: WorkerStatus;
  role: WorkerRole;
  cpuUsage: number;
  memoryUsage: string;
  checkpointId?: string;
  lastHeartbeat?: string;
  createdAt?: string;
  updatedAt?: string;
}

export type WorkflowStatus =
  | "RUNNING"
  | "FAILED"
  | "COMPLETED"
  | "PAUSED";

export interface Workflow {
  id: string;
  name: string;
  status: WorkflowStatus;
  createdAt?: string;
  updatedAt?: string;
  workerIds?: string[];
}

export interface Checkpoint {
  id: string;
  createdAt: string;
  size: string;
  storagePath?: string;
  checksum?: string;
  version?: number;
}

export type TelemetryLevel =
  | "INFO"
  | "WARN"
  | "ERROR"
  | "SUCCESS";

export interface TelemetryEvent {
  level: TelemetryLevel;
  message: string;
  timestamp: string;
  workerId?: string;
  workflowId?: string;
  metadata?: Record<string, unknown>;
}

export type VariableDiffType =
  | "modified"
  | "added"
  | "removed"
  | "unchanged";

export interface VariableDiffData {
  key: string;
  before?: string;
  after?: string;
  type: VariableDiffType;
}

export type VariableDiff = VariableDiffData;

export interface ApiResponse<T> {
  data: T;
  message?: string;
  timestamp: string;
}

export interface Pagination {
  page: number;
  limit: number;
  total: number;
}

export interface HealthStatus {
  status:
    | "HEALTHY"
    | "DEGRADED"
    | "UNHEALTHY";
  uptime: number;
  version: string;
}
