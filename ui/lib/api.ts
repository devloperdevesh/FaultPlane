import type { DashboardMetrics } from "./types";

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080";

function joinUrl(base: string, endpoint: string): string {
  return `${base.replace(/\/+$/, "")}/${endpoint.replace(/^\/+/, "")}`;
}

async function request<T>(endpoint: string): Promise<T> {
  const response = await fetch(joinUrl(API_URL, endpoint), {
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error(
      `FaultPlane API request failed: ${response.status} ${response.statusText}`,
    );
  }

  return response.json() as Promise<T>;
}

export async function getMetrics(): Promise<DashboardMetrics> {
  return request<DashboardMetrics>("/api/metrics");
}

export interface RuntimeWorker {
  id: string;
  status: string;
  cpu: number;
  memory: number;
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
