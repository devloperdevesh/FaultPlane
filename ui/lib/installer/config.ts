/**
 * Configuration model for the Kubernetes sidecar manifest generator.
 *
 * Environment variable names mirror what the daemon reads in
 * `internal/config/config.go` (FAULTPLANE_HOST, FAULTPLANE_PORT,
 * FAULTPLANE_LOG_LEVEL). The optional upstream addresses mirror the
 * repository's docker-compose example.
 */

export const LOG_LEVELS = ["debug", "info", "warn", "error"] as const;

export type LogLevel = (typeof LOG_LEVELS)[number];

export interface InstallerConfig {
  namespace: string;
  workloadName: string;
  replicas: number;
  imageRepository: string;
  imageTag: string;
  containerPort: number;
  logLevel: LogLevel;
  primaryAddr: string;
  fallbackAddr: string;
  cpuRequest: string;
  memoryRequest: string;
  cpuLimit: string;
  memoryLimit: string;
  enableProbes: boolean;
  enablePrometheusAnnotations: boolean;
}

/** Health endpoint served by the gateway (see deployments/daemon/healthcheck). */
export const HEALTH_PATH = "/health";

/** Prometheus endpoint served by the gateway (see internal/gateway/router.go). */
export const METRICS_PATH = "/api/v1/metrics/prometheus";

export const DEFAULT_CONFIG: InstallerConfig = {
  namespace: "default",
  workloadName: "faultplane",
  replicas: 1,
  // Editable placeholder. Replace with the image reference and tag you publish.
  imageRepository: "ghcr.io/devloperdevesh/faultplane",
  imageTag: "latest",
  containerPort: 8080,
  logLevel: "info",
  primaryAddr: "",
  fallbackAddr: "",
  cpuRequest: "100m",
  memoryRequest: "128Mi",
  cpuLimit: "500m",
  memoryLimit: "512Mi",
  enableProbes: true,
  enablePrometheusAnnotations: false,
};
