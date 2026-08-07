// ===============================
// Runtime Worker Contract
// ===============================

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


// ===============================
// Workflow Contract
// ===============================

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


// ===============================
// Checkpoint Contract
// ===============================

export interface Checkpoint {

  id: string;

  createdAt: string;

  size: string;

  storagePath?: string;

  checksum?: string;

  version?: number;

}


// ===============================
// Metrics Contract
// ===============================

export interface Metric {

  name: string;

  value: number;

  unit?: string;

  timestamp?: string;

}


// ===============================
// Telemetry Contract
// ===============================

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


// ===============================
// GitOps State Diff Contract
// ===============================

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


// Existing components compatibility
export type VariableDiff = VariableDiffData;


// ===============================
// API Response Wrapper
// ===============================

export interface ApiResponse<T> {

  data: T;

  message?: string;

  timestamp: string;

}


// ===============================
// Pagination Contract
// ===============================

export interface Pagination {

  page: number;

  limit: number;

  total: number;

}


// ===============================
// Health Status
// ===============================

export interface HealthStatus {

  status:
    | "HEALTHY"
    | "DEGRADED"
    | "UNHEALTHY";

  uptime: number;

  version: string;

}