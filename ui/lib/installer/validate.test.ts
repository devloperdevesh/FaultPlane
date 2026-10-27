import { describe, expect, it } from "vitest";

import { DEFAULT_CONFIG, type InstallerConfig } from "./config";
import {
  isDnsLabel,
  isHostPort,
  parseCpuMillicores,
  parseMemoryBytes,
  validateConfig,
} from "./validate";

const withConfig = (overrides: Partial<InstallerConfig>): InstallerConfig => ({
  ...DEFAULT_CONFIG,
  ...overrides,
});

describe("validateConfig", () => {
  it("accepts the default configuration", () => {
    const result = validateConfig(DEFAULT_CONFIG);
    expect(result.valid).toBe(true);
    expect(result.errors).toEqual({});
  });

  it("warns about the mutable latest tag without blocking", () => {
    const result = validateConfig(DEFAULT_CONFIG);
    expect(result.warnings.join(" ")).toContain("latest");
    expect(validateConfig(withConfig({ imageTag: "v1.2.3" })).warnings).toEqual([]);
  });

  it.each([
    ["Upper", "uppercase letters"],
    ["has_underscore", "underscores"],
    ["-leading", "a leading dash"],
    ["trailing-", "a trailing dash"],
    ["", "an empty value"],
    ["a".repeat(64), "more than 63 characters"],
  ])("rejects namespace %j (%s)", (namespace) => {
    expect(validateConfig(withConfig({ namespace })).errors.namespace).toBeDefined();
  });

  it("rejects an invalid workload name", () => {
    expect(validateConfig(withConfig({ workloadName: "My App" })).errors.workloadName).toBeDefined();
  });

  it.each([0, 51, 1.5, Number.NaN])("rejects replicas %s", (replicas) => {
    expect(validateConfig(withConfig({ replicas })).errors.replicas).toBeDefined();
  });

  it.each([
    ["ghcr.io/org/faultplane:latest", "a tag in the repository"],
    ["Ghcr.io/Org/Name", "uppercase letters"],
    ["has space/name", "spaces"],
    ["", "an empty value"],
  ])("rejects image repository %j (%s)", (imageRepository) => {
    expect(validateConfig(withConfig({ imageRepository })).errors.imageRepository).toBeDefined();
  });

  it("accepts registries with ports and nested paths", () => {
    const result = validateConfig(withConfig({ imageRepository: "registry.local:5000/team/faultplane" }));
    expect(result.errors.imageRepository).toBeUndefined();
  });

  it("rejects an invalid image tag", () => {
    expect(validateConfig(withConfig({ imageTag: "bad tag!" })).errors.imageTag).toBeDefined();
  });

  it.each([80, 1023, 65536, 0, 8080.5, Number.NaN])("rejects container port %s", (containerPort) => {
    expect(validateConfig(withConfig({ containerPort })).errors.containerPort).toBeDefined();
  });

  it("accepts boundary ports 1024 and 65535", () => {
    expect(validateConfig(withConfig({ containerPort: 1024 })).errors.containerPort).toBeUndefined();
    expect(validateConfig(withConfig({ containerPort: 65535 })).errors.containerPort).toBeUndefined();
  });

  it("rejects an unknown log level", () => {
    const config = withConfig({ logLevel: "trace" as InstallerConfig["logLevel"] });
    expect(validateConfig(config).errors.logLevel).toBeDefined();
  });

  it("validates upstream addresses as host:port", () => {
    const bad = validateConfig(withConfig({ primaryAddr: "no-port", fallbackAddr: "host:99999" }));
    expect(bad.errors.primaryAddr).toBeDefined();
    expect(bad.errors.fallbackAddr).toBeDefined();

    const good = validateConfig(
      withConfig({ primaryAddr: "backend.default.svc:8080", fallbackAddr: "10.0.0.5:9000" }),
    );
    expect(good.valid).toBe(true);
  });

  it("requires a primary upstream when a fallback is set", () => {
    const result = validateConfig(withConfig({ primaryAddr: "", fallbackAddr: "backup:8080" }));
    expect(result.errors.fallbackAddr).toContain("primary");
  });

  it("rejects malformed resource quantities", () => {
    const result = validateConfig(
      withConfig({ cpuRequest: "abc", cpuLimit: "0", memoryRequest: "12XB", memoryLimit: "" }),
    );
    expect(result.errors.cpuRequest).toBeDefined();
    expect(result.errors.cpuLimit).toBeDefined();
    expect(result.errors.memoryRequest).toBeDefined();
    expect(result.errors.memoryLimit).toBeDefined();
  });

  it("rejects requests that exceed limits", () => {
    const cpu = validateConfig(withConfig({ cpuRequest: "1", cpuLimit: "500m" }));
    expect(cpu.errors.cpuRequest).toContain("limit");

    const memory = validateConfig(withConfig({ memoryRequest: "1Gi", memoryLimit: "512Mi" }));
    expect(memory.errors.memoryRequest).toContain("limit");
  });

  it("accepts requests equal to limits", () => {
    const result = validateConfig(withConfig({ cpuRequest: "500m", cpuLimit: "0.5" }));
    expect(result.errors.cpuRequest).toBeUndefined();
  });
});

describe("quantity helpers", () => {
  it("parses CPU quantities to millicores", () => {
    expect(parseCpuMillicores("250m")).toBe(250);
    expect(parseCpuMillicores("1")).toBe(1000);
    expect(parseCpuMillicores("0.5")).toBe(500);
    expect(parseCpuMillicores("0.0001")).toBeNull();
    expect(parseCpuMillicores("-1")).toBeNull();
  });

  it("parses memory quantities to bytes", () => {
    expect(parseMemoryBytes("1Ki")).toBe(1024);
    expect(parseMemoryBytes("128Mi")).toBe(128 * 1024 * 1024);
    expect(parseMemoryBytes("1G")).toBe(1e9);
    expect(parseMemoryBytes("0")).toBeNull();
    expect(parseMemoryBytes("1mi")).toBeNull();
  });

  it("recognises DNS labels and host:port values", () => {
    expect(isDnsLabel("faultplane-1")).toBe(true);
    expect(isDnsLabel("Fault")).toBe(false);
    expect(isHostPort("svc.ns:80")).toBe(true);
    expect(isHostPort("svc.ns")).toBe(false);
  });
});
