import type {
  DashboardMetrics,
  RuntimeWorker,
  TelemetryResponse,
  LogsResponse,
  NetworkEventsResponse,
  TopologySnapshot,
  HealthResponse,
  EbpfStatusResponse,
  EbpfHooksResponse,
  EbpfEventsResponse,
} from "./types";

const API_URL = (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080").replace(/\/+$/, "");

function joinUrl(endpoint: string): string {
  return `${API_URL}/${endpoint.replace(/^\/+/, "")}`;
}

async function request<T>(endpoint: string, init?: RequestInit): Promise<T> {
  const response = await fetch(joinUrl(endpoint), {
    ...init,
    cache: "no-store",
    headers: {
      Accept: "application/json",
      ...(init?.headers ?? {}),
    },
  });

  if (!response.ok) {
    throw new Error(
      `FaultPlane API request failed: ${response.status} ${response.statusText}`,
    );
  }

  return response.json() as Promise<T>;
}

export function getApiUrl(): string {
  return API_URL;
}

export async function getHealth(): Promise<HealthResponse> {
  return request<HealthResponse>("/health");
}

export async function getMetrics(): Promise<DashboardMetrics> {
  return request<DashboardMetrics>("/api/metrics");
}

export async function getWorkers(): Promise<RuntimeWorker[]> {
  return request<RuntimeWorker[]>("/api/workers");
}

export async function getTelemetry(): Promise<TelemetryResponse> {
  return request<TelemetryResponse>("/api/telemetry");
}

export async function getLogs(): Promise<LogsResponse> {
  return request<LogsResponse>("/api/logs");
}

export async function getNetworkEvents(): Promise<NetworkEventsResponse> {
  return request<NetworkEventsResponse>("/api/network/events");
}

export async function getTopology(): Promise<TopologySnapshot> {
  return request<TopologySnapshot>("/api/topology");
}

export async function getEbpfStatus(): Promise<EbpfStatusResponse> {
  return request<EbpfStatusResponse>("/api/ebpf/status");
}

export async function getEbpfEvents(): Promise<EbpfEventsResponse> {
  return request<EbpfEventsResponse>("/api/ebpf/events");
}

export async function getEbpfHooks(): Promise<EbpfHooksResponse> {
  return request<EbpfHooksResponse>("/api/ebpf/hooks");
}

export type { DashboardMetrics, RuntimeWorker, TelemetryEvent, TelemetryResponse, LogsResponse, NetworkEventsResponse, TopologySnapshot, TopologyNode, TopologyConnection, HealthResponse, EbpfStatusResponse, EbpfHook, EbpfHooksResponse, EbpfEvent, EbpfEventsResponse } from './types';