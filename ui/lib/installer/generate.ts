import type { InstallerConfig } from "./config";
import { HEALTH_PATH, METRICS_PATH } from "./config";
import type { ValidationErrors } from "./validate";
import { validateConfig } from "./validate";

export interface GeneratedManifests {
  /** Container entry to paste under spec.template.spec.containers of an existing workload. */
  sidecar: string;
  /** Standalone Deployment plus ClusterIP Service. */
  deployment: string;
  /** values.yaml for a Helm chart you maintain (this repository ships none yet). */
  helmValues: string;
}

export type GenerateResult =
  | { ok: true; manifests: GeneratedManifests; warnings: string[] }
  | { ok: false; errors: ValidationErrors };

export const FILE_NAMES: Record<keyof GeneratedManifests, string> = {
  sidecar: "faultplane-sidecar-container.yaml",
  deployment: "faultplane-deployment.yaml",
  helmValues: "faultplane-values.yaml",
};

/**
 * Every string scalar is emitted as a JSON string, which is a valid YAML
 * double-quoted scalar. This keeps values like "true" or "123" as strings and
 * makes injection through free-text fields impossible.
 */
const q = (value: string | number): string => JSON.stringify(String(value));

function indent(lines: string[], spaces: number): string[] {
  const pad = " ".repeat(spaces);
  return lines.map((line) => (line === "" ? line : pad + line));
}

function containerLines(config: InstallerConfig, name: string): string[] {
  const lines: string[] = [
    `- name: ${q(name)}`,
    `  image: ${q(`${config.imageRepository}:${config.imageTag}`)}`,
    "  imagePullPolicy: IfNotPresent",
    "  ports:",
    "    - name: faultplane",
    `      containerPort: ${config.containerPort}`,
    "      protocol: TCP",
    "  env:",
    "    - name: FAULTPLANE_HOST",
    '      value: "0.0.0.0"',
    "    - name: FAULTPLANE_PORT",
    `      value: ${q(config.containerPort)}`,
    "    - name: FAULTPLANE_LOG_LEVEL",
    `      value: ${q(config.logLevel)}`,
  ];

  if (config.primaryAddr) {
    lines.push("    - name: FAULTPLANE_PRIMARY_ADDR", `      value: ${q(config.primaryAddr)}`);
  }
  if (config.fallbackAddr) {
    lines.push("    - name: FAULTPLANE_FALLBACK_ADDR", `      value: ${q(config.fallbackAddr)}`);
  }

  lines.push(
    "  resources:",
    "    requests:",
    `      cpu: ${q(config.cpuRequest)}`,
    `      memory: ${q(config.memoryRequest)}`,
    "    limits:",
    `      cpu: ${q(config.cpuLimit)}`,
    `      memory: ${q(config.memoryLimit)}`,
  );

  if (config.enableProbes) {
    lines.push(
      "  readinessProbe:",
      "    httpGet:",
      `      path: ${HEALTH_PATH}`,
      "      port: faultplane",
      "    initialDelaySeconds: 3",
      "    periodSeconds: 10",
      "  livenessProbe:",
      "    httpGet:",
      `      path: ${HEALTH_PATH}`,
      "      port: faultplane",
      "    initialDelaySeconds: 10",
      "    periodSeconds: 10",
    );
  }

  // Deliberately unprivileged: no privileged mode and no added capabilities.
  lines.push(
    "  securityContext:",
    "    allowPrivilegeEscalation: false",
    "    runAsNonRoot: true",
    "    capabilities:",
    "      drop:",
    "        - ALL",
  );

  return lines;
}

function buildSidecar(config: InstallerConfig): string {
  const header = [
    "# Paste this entry under spec.template.spec.containers of an existing workload.",
    "# Preview output: review it before applying. It does not grant privileged or eBPF access.",
  ];
  return [...header, ...containerLines(config, "faultplane-sidecar")].join("\n") + "\n";
}

function buildDeployment(config: InstallerConfig): string {
  const name = config.workloadName;
  const labels = [
    `app.kubernetes.io/name: ${q(name)}`,
    'app.kubernetes.io/component: "proxy"',
  ];

  const podAnnotations = config.enablePrometheusAnnotations
    ? [
        "annotations:",
        '  prometheus.io/scrape: "true"',
        `  prometheus.io/port: ${q(config.containerPort)}`,
        `  prometheus.io/path: ${q(METRICS_PATH)}`,
      ]
    : [];

  const lines: string[] = [
    "# Preview output: review it before applying. It does not grant privileged or eBPF access.",
    "apiVersion: apps/v1",
    "kind: Deployment",
    "metadata:",
    `  name: ${q(name)}`,
    `  namespace: ${q(config.namespace)}`,
    "  labels:",
    ...indent(labels, 4),
    "spec:",
    `  replicas: ${config.replicas}`,
    "  selector:",
    "    matchLabels:",
    `      app.kubernetes.io/name: ${q(name)}`,
    "  template:",
    "    metadata:",
    "      labels:",
    ...indent(labels, 8),
    ...indent(podAnnotations, 6),
    "    spec:",
    "      containers:",
    ...indent(containerLines(config, name), 8),
    "---",
    "apiVersion: v1",
    "kind: Service",
    "metadata:",
    `  name: ${q(name)}`,
    `  namespace: ${q(config.namespace)}`,
    "  labels:",
    ...indent(labels, 4),
    "spec:",
    "  type: ClusterIP",
    "  selector:",
    `    app.kubernetes.io/name: ${q(name)}`,
    "  ports:",
    "    - name: faultplane",
    `      port: ${config.containerPort}`,
    "      targetPort: faultplane",
    "      protocol: TCP",
  ];

  return lines.join("\n") + "\n";
}

function buildHelmValues(config: InstallerConfig): string {
  const lines: string[] = [
    "# values.yaml for a FaultPlane Helm chart that you maintain.",
    "# This repository does not ship a chart yet. Preview output: review before use.",
    `replicaCount: ${config.replicas}`,
    `nameOverride: ${q(config.workloadName)}`,
    `namespace: ${q(config.namespace)}`,
    "image:",
    `  repository: ${q(config.imageRepository)}`,
    `  tag: ${q(config.imageTag)}`,
    "  pullPolicy: IfNotPresent",
    "service:",
    "  type: ClusterIP",
    `  port: ${config.containerPort}`,
    "faultplane:",
    '  host: "0.0.0.0"',
    `  port: ${config.containerPort}`,
    `  logLevel: ${q(config.logLevel)}`,
  ];

  if (config.primaryAddr) lines.push(`  primaryAddr: ${q(config.primaryAddr)}`);
  if (config.fallbackAddr) lines.push(`  fallbackAddr: ${q(config.fallbackAddr)}`);

  lines.push(
    "resources:",
    "  requests:",
    `    cpu: ${q(config.cpuRequest)}`,
    `    memory: ${q(config.memoryRequest)}`,
    "  limits:",
    `    cpu: ${q(config.cpuLimit)}`,
    `    memory: ${q(config.memoryLimit)}`,
    "probes:",
    `  enabled: ${config.enableProbes}`,
    `  path: ${HEALTH_PATH}`,
    "prometheus:",
    `  scrape: ${config.enablePrometheusAnnotations}`,
    `  path: ${q(METRICS_PATH)}`,
    "securityContext:",
    "  allowPrivilegeEscalation: false",
    "  runAsNonRoot: true",
    "  capabilities:",
    "    drop:",
    "      - ALL",
  );

  return lines.join("\n") + "\n";
}

/** Validates the configuration and, if valid, renders all manifests. */
export function generateManifests(config: InstallerConfig): GenerateResult {
  const validation = validateConfig(config);

  if (!validation.valid) {
    return { ok: false, errors: validation.errors };
  }

  return {
    ok: true,
    warnings: validation.warnings,
    manifests: {
      sidecar: buildSidecar(config),
      deployment: buildDeployment(config),
      helmValues: buildHelmValues(config),
    },
  };
}
