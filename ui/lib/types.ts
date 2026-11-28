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
  dma_bound?: boolean;
  expert_id?: number;
  dma_latency_ms?: number;
}

export interface DMAChannelSnapshot {
  channel_id: number;
  gpu_bus_id: string;
  bandwidth_gbs: number;
  active_streams: number;
  total_transfers: number;
  average_latency_ms: number;
  sub_two_ms_bound_met: boolean;
  active_expert_id: number;
}

export interface DPDKPortSnapshot {
  port_id: number;
  device_pci_addr: string;
  link_speed_mbps: number;
  poll_mode_driver: string;
  kernel_bypass_rate: number;
  zero_copy_hot_swaps: number;
  context_cache_bits_routed: number;
  hardware_ring_active: boolean;
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
  pcie_dma_bound?: boolean;
  dpdk_bypass_active?: boolean;
  dma_throughput_gbs?: number;
  average_dma_latency_ms?: number;
  sub_two_ms_congruent?: boolean;
  moe_experts_offloaded?: number;
  os_kernel_bypass_percent?: number;
  dma_channels?: DMAChannelSnapshot[];
  dpdk_port?: DPDKPortSnapshot;
  slots: SlotSnapshot[];
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
export interface HardwareProfileResponse {
  has_avx512f: boolean;
  has_avx512bw: boolean;
  has_avx512cd: boolean;
  has_avx512dq: boolean;
  has_avx512vl: boolean;
  cgo_available: boolean;
  active_track: string;
  vector_bits: number;
  parallel_lanes: number;
}