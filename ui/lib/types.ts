export type VariableDiffType = "modified" | "added" | "removed" | "unchanged";

export interface VariableDiffData {
  key: string;
  before: string;
  after: string;
  type: VariableDiffType;
}

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

export interface RuntimeWorker {
  id: string;
  region?: string;
  status: string;
  cpu: number;
  memory: number;
  role?: string;
  checkpointId?: string;
  lastHeartbeat?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface TelemetryEvent {
  type: string;
  timestamp: string;
  value?: number;
  metadata?: Record<string, string>;
}

export interface TelemetryResponse {
  events: TelemetryEvent[];
}

export interface LogsResponse {
  logs: TelemetryEvent[];
}

export interface NetworkEventsResponse {
  events: TelemetryEvent[];
}

export interface TopologyNode {
  id: string;
  region?: string;
  type: string;
  name: string;
  status: string;
}

export interface TopologyConnection {
  id: string;
  source: string;
  target: string;
  type: string;
  status: string;
}

export interface TopologySnapshot {
  nodes: TopologyNode[];
  connections: TopologyConnection[];
  updated_at: string;
}

export interface HealthResponse {
  status: string;
  service: string;
  timestamp?: string;
}

export interface EbpfStatusResponse {
  status: string;
}

export interface EbpfHook {
  name: string;
  status?: string;
  [key: string]: unknown;
}

export interface EbpfHooksResponse {
  hooks: EbpfHook[];
}

export interface EbpfEvent {
  [key: string]: unknown;
}

export interface EbpfEventsResponse {
  events: EbpfEvent[];
}

export interface ComplianceMetricSummary {
  uptime_percentage: number;
  total_monitored_hours: number;
  total_inference_loops: number;
  context_loss_events: number;
  failover_events: number;
  average_failover_latency_ms: number;
  data_isolation_enforced: boolean;
  silicon_vault_active: boolean;
  cache_alignment_bytes: number;
  compliance_standards_passed: string[];
}

export interface ComplianceAuditLog {
  id: string;
  timestamp: string;
  timeframe: string;
  tenant_id: string;
  uptime_percentage: number;
  failover_latency_ms: number;
  context_loss_count: number;
  isolation_status: string;
  cache_aligned: boolean;
  standard: string;
  signature_hash: string;
}

export interface ComplianceLedgerResponse {
  ledger_id: string;
  generated_at: string;
  signer_identity: string;
  verification_key_hash: string;
  ledger_signature: string;
  summary: ComplianceMetricSummary;
  logs: ComplianceAuditLog[];
}
