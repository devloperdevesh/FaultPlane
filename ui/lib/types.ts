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

export interface SlotSnapshot {
  index: number;
  slot_id: number;
  tenant_id: string;
  domain_key: number;
  status: "ACTIVE" | "FENCED" | "FREE" | string;
  active_tokens: number;
  access_count: number;
  memory_offset_hex: string;
  cache_line_aligned: boolean;
  color_channel: string;
}

export interface BufferPoolSnapshot {
  total_slots: number;
  active_slots: number;
  fenced_slots: number;
  cache_line_size: number;
  cache_alignment_congruent: boolean;
  l1_l2_thrashing_rate: number;
  gc_pauses_ms: number;
  active_tenants_count: number;
  slots: SlotSnapshot[];
}
