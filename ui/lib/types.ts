// Runtime Worker Contract

export interface Worker {
  id: string;

  status:
    | "ACTIVE"
    | "FAILED"
    | "RECOVERING";

  cpu: number;

  memory: string;

  checkpoint: string;

  role:
    | "PRIMARY"
    | "STANDBY";
}


// Workflow Contract

export interface Workflow {

  id: string;

  name: string;

  status:
    | "RUNNING"
    | "FAILED";

}


// Checkpoint Contract

export interface Checkpoint {

  id: string;

  createdAt: string;

  size: string;

}


// Metrics Contract

export interface Metric {

  name: string;

  value: string;

}


// Telemetry Contract

export interface TelemetryEvent {

  level:
    | "INFO"
    | "WARN"
    | "ERROR"
    | "SUCCESS";

  message: string;

  timestamp: string;

}


// GitOps State Diff Contract

export type VariableDiffType =
  | "modified"
  | "added"
  | "removed"
  | "unchanged";


export interface VariableDiffData {

  key: string;

  before: string;

  after: string;

  type: VariableDiffType;

}