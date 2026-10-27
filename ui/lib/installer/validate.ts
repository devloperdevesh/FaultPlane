import type { InstallerConfig } from "./config";
import { LOG_LEVELS } from "./config";

export type FieldName = keyof InstallerConfig;

export type ValidationErrors = Partial<Record<FieldName, string>>;

export interface ValidationResult {
  valid: boolean;
  errors: ValidationErrors;
  warnings: string[];
}

const DNS_LABEL = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/;

// Optional registry host (with optional port) followed by lowercase path parts.
const IMAGE_REPOSITORY =
  /^(?:[a-z0-9]+(?:[._-][a-z0-9]+)*(?::[0-9]+)?\/)?[a-z0-9]+(?:[._-][a-z0-9]+)*(?:\/[a-z0-9]+(?:[._-][a-z0-9]+)*)*$/;

const IMAGE_TAG = /^[A-Za-z0-9_][A-Za-z0-9_.-]{0,127}$/;

const HOST_PORT = /^([A-Za-z0-9](?:[-A-Za-z0-9.]*[A-Za-z0-9])?):(\d{1,5})$/;

const MEMORY_UNITS: Record<string, number> = {
  "": 1,
  k: 1e3,
  M: 1e6,
  G: 1e9,
  T: 1e12,
  P: 1e15,
  E: 1e18,
  Ki: 1024,
  Mi: 1024 ** 2,
  Gi: 1024 ** 3,
  Ti: 1024 ** 4,
  Pi: 1024 ** 5,
  Ei: 1024 ** 6,
};

export function isDnsLabel(value: string): boolean {
  return value.length >= 1 && value.length <= 63 && DNS_LABEL.test(value);
}

export function isPort(value: number): boolean {
  return Number.isInteger(value) && value >= 1 && value <= 65535;
}

export function isHostPort(value: string): boolean {
  const match = HOST_PORT.exec(value);
  if (!match) return false;
  if (match[1].length > 253) return false;
  return isPort(Number(match[2]));
}

/** Parses a Kubernetes CPU quantity ("250m", "1", "0.5") to millicores. */
export function parseCpuMillicores(value: string): number | null {
  const milli = /^(\d+)m$/.exec(value);
  if (milli) {
    const n = Number(milli[1]);
    return n > 0 ? n : null;
  }
  if (/^\d+(\.\d+)?$/.test(value)) {
    const millis = Number(value) * 1000;
    if (!(millis > 0)) return null;
    if (Math.abs(millis - Math.round(millis)) > 1e-6) return null;
    return Math.round(millis);
  }
  return null;
}

/** Parses a Kubernetes memory quantity ("128Mi", "1G", "1048576") to bytes. */
export function parseMemoryBytes(value: string): number | null {
  const match = /^(\d+(?:\.\d+)?)(Ki|Mi|Gi|Ti|Pi|Ei|k|M|G|T|P|E)?$/.exec(value);
  if (!match) return null;
  const bytes = Number(match[1]) * MEMORY_UNITS[match[2] ?? ""];
  return bytes > 0 ? bytes : null;
}

const DNS_LABEL_MESSAGE =
  "Use 1-63 lowercase letters, digits or '-', starting and ending with a letter or digit.";

export function validateConfig(config: InstallerConfig): ValidationResult {
  const errors: ValidationErrors = {};
  const warnings: string[] = [];

  // Identity
  if (!config.namespace) {
    errors.namespace = "Namespace is required.";
  } else if (!isDnsLabel(config.namespace)) {
    errors.namespace = `Invalid namespace. ${DNS_LABEL_MESSAGE}`;
  }

  if (!config.workloadName) {
    errors.workloadName = "Workload name is required.";
  } else if (!isDnsLabel(config.workloadName)) {
    errors.workloadName = `Invalid workload name. ${DNS_LABEL_MESSAGE}`;
  }

  if (!Number.isInteger(config.replicas) || config.replicas < 1 || config.replicas > 50) {
    errors.replicas = "Replicas must be a whole number between 1 and 50.";
  }

  // Image
  if (!config.imageRepository) {
    errors.imageRepository = "Image repository is required.";
  } else if (config.imageRepository.length > 255 || !IMAGE_REPOSITORY.test(config.imageRepository)) {
    errors.imageRepository =
      "Use a lowercase repository such as ghcr.io/org/name, without a tag or spaces.";
  }

  if (!config.imageTag) {
    errors.imageTag = "Image tag is required.";
  } else if (!IMAGE_TAG.test(config.imageTag)) {
    errors.imageTag =
      "Tags may contain letters, digits, '_', '.' and '-' (max 128 characters).";
  } else if (config.imageTag === "latest") {
    warnings.push(
      "The tag 'latest' is not reproducible. Pin a published version tag for production.",
    );
  }

  // Networking
  if (!Number.isInteger(config.containerPort)) {
    errors.containerPort = "Container port must be a whole number.";
  } else if (config.containerPort < 1024) {
    errors.containerPort =
      "Ports below 1024 need extra privileges. The generated container runs as non-root, so use 1024-65535.";
  } else if (config.containerPort > 65535) {
    errors.containerPort = "Container port must be 65535 or lower.";
  }

  if (!LOG_LEVELS.includes(config.logLevel)) {
    errors.logLevel = `Log level must be one of: ${LOG_LEVELS.join(", ")}.`;
  }

  if (config.primaryAddr && !isHostPort(config.primaryAddr)) {
    errors.primaryAddr = "Use host:port, for example backend.default.svc:8080.";
  }
  if (config.fallbackAddr && !isHostPort(config.fallbackAddr)) {
    errors.fallbackAddr = "Use host:port, for example backup.default.svc:8080.";
  }
  if (config.fallbackAddr && !config.primaryAddr && !errors.fallbackAddr) {
    errors.fallbackAddr = "Set a primary upstream before adding a fallback.";
  }

  // Resources
  const cpuRequest = parseCpuMillicores(config.cpuRequest);
  const cpuLimit = parseCpuMillicores(config.cpuLimit);
  const memoryRequest = parseMemoryBytes(config.memoryRequest);
  const memoryLimit = parseMemoryBytes(config.memoryLimit);

  if (cpuRequest === null) {
    errors.cpuRequest = "Use a CPU quantity such as 100m or 0.5.";
  }
  if (cpuLimit === null) {
    errors.cpuLimit = "Use a CPU quantity such as 500m or 1.";
  }
  if (memoryRequest === null) {
    errors.memoryRequest = "Use a memory quantity such as 128Mi or 1Gi.";
  }
  if (memoryLimit === null) {
    errors.memoryLimit = "Use a memory quantity such as 512Mi or 1Gi.";
  }
  if (cpuRequest !== null && cpuLimit !== null && cpuRequest > cpuLimit) {
    errors.cpuRequest = "CPU request must not exceed the CPU limit.";
  }
  if (memoryRequest !== null && memoryLimit !== null && memoryRequest > memoryLimit) {
    errors.memoryRequest = "Memory request must not exceed the memory limit.";
  }

  return { valid: Object.keys(errors).length === 0, errors, warnings };
}
