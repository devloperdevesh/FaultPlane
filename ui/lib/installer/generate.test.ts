import { parse, parseAllDocuments } from "yaml";
import { describe, expect, it } from "vitest";

import { DEFAULT_CONFIG, HEALTH_PATH, METRICS_PATH, type InstallerConfig } from "./config";
import { generateManifests, type GeneratedManifests } from "./generate";

const withConfig = (overrides: Partial<InstallerConfig>): InstallerConfig => ({
  ...DEFAULT_CONFIG,
  ...overrides,
});

function manifestsFor(overrides: Partial<InstallerConfig> = {}): GeneratedManifests {
  const result = generateManifests(withConfig(overrides));
  if (!result.ok) {
    throw new Error(`expected valid config, got ${JSON.stringify(result.errors)}`);
  }
  return result.manifests;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Yaml = any;

describe("generateManifests", () => {
  it("returns errors and no manifests for an invalid configuration", () => {
    const result = generateManifests(withConfig({ namespace: "Bad Name", containerPort: 80 }));
    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.errors.namespace).toBeDefined();
      expect(result.errors.containerPort).toBeDefined();
    }
  });

  it("produces YAML that parses without errors", () => {
    const manifests = manifestsFor({
      primaryAddr: "backend.default.svc:8080",
      fallbackAddr: "backup.default.svc:8080",
      enablePrometheusAnnotations: true,
    });

    for (const text of [manifests.sidecar, manifests.deployment, manifests.helmValues]) {
      for (const doc of parseAllDocuments(text)) {
        expect(doc.errors).toEqual([]);
      }
    }
  });

  describe("sidecar container", () => {
    it("reflects the selected configuration", () => {
      const [container] = parse(
        manifestsFor({
          imageRepository: "registry.local:5000/team/faultplane",
          imageTag: "v1.4.0",
          containerPort: 9090,
          logLevel: "debug",
          cpuRequest: "250m",
          memoryLimit: "1Gi",
        }).sidecar,
      ) as Yaml[];

      expect(container.name).toBe("faultplane-sidecar");
      expect(container.image).toBe("registry.local:5000/team/faultplane:v1.4.0");
      expect(container.ports[0].containerPort).toBe(9090);
      expect(container.resources.requests.cpu).toBe("250m");
      expect(container.resources.limits.memory).toBe("1Gi");

      const env = Object.fromEntries(container.env.map((e: Yaml) => [e.name, e.value]));
      expect(env.FAULTPLANE_PORT).toBe("9090");
      expect(env.FAULTPLANE_LOG_LEVEL).toBe("debug");
      expect(env.FAULTPLANE_HOST).toBe("0.0.0.0");
    });

    it("emits upstream addresses only when set", () => {
      const without = parse(manifestsFor().sidecar) as Yaml[];
      expect(JSON.stringify(without[0].env)).not.toContain("FAULTPLANE_PRIMARY_ADDR");

      const withAddrs = parse(
        manifestsFor({ primaryAddr: "a.b:1111", fallbackAddr: "c.d:2222" }).sidecar,
      ) as Yaml[];
      const env = Object.fromEntries(withAddrs[0].env.map((e: Yaml) => [e.name, e.value]));
      expect(env.FAULTPLANE_PRIMARY_ADDR).toBe("a.b:1111");
      expect(env.FAULTPLANE_FALLBACK_ADDR).toBe("c.d:2222");
    });

    it("includes health probes only when enabled", () => {
      const on = (parse(manifestsFor({ enableProbes: true }).sidecar) as Yaml[])[0];
      expect(on.readinessProbe.httpGet.path).toBe(HEALTH_PATH);
      expect(on.livenessProbe.httpGet.port).toBe("faultplane");

      const off = (parse(manifestsFor({ enableProbes: false }).sidecar) as Yaml[])[0];
      expect(off.readinessProbe).toBeUndefined();
      expect(off.livenessProbe).toBeUndefined();
    });

    it("never requests privileged mode or added capabilities", () => {
      const output = manifestsFor().sidecar + manifestsFor().deployment + manifestsFor().helmValues;
      expect(output).not.toMatch(/privileged:\s*true/);
      expect(output).not.toMatch(/\badd:/);
      expect(output).not.toMatch(/hostNetwork/);

      const [container] = parse(manifestsFor().sidecar) as Yaml[];
      expect(container.securityContext.allowPrivilegeEscalation).toBe(false);
      expect(container.securityContext.runAsNonRoot).toBe(true);
      expect(container.securityContext.capabilities.drop).toEqual(["ALL"]);
    });
  });

  describe("deployment and service", () => {
    it("renders a Deployment whose selector matches its pod labels", () => {
      const docs = parseAllDocuments(
        manifestsFor({ namespace: "payments", workloadName: "fp-gateway", replicas: 3 }).deployment,
      ).map((d) => d.toJS() as Yaml);

      const [deployment, service] = docs;
      expect(deployment.kind).toBe("Deployment");
      expect(deployment.metadata.namespace).toBe("payments");
      expect(deployment.metadata.name).toBe("fp-gateway");
      expect(deployment.spec.replicas).toBe(3);
      expect(deployment.spec.selector.matchLabels).toEqual({ "app.kubernetes.io/name": "fp-gateway" });
      expect(deployment.spec.template.metadata.labels["app.kubernetes.io/name"]).toBe("fp-gateway");
      expect(deployment.spec.template.spec.containers[0].name).toBe("fp-gateway");

      expect(service.kind).toBe("Service");
      expect(service.metadata.namespace).toBe("payments");
      expect(service.spec.selector).toEqual({ "app.kubernetes.io/name": "fp-gateway" });
      expect(service.spec.ports[0].targetPort).toBe("faultplane");
    });

    it("adds Prometheus annotations only when requested", () => {
      const off = parseAllDocuments(manifestsFor().deployment)[0].toJS() as Yaml;
      expect(off.spec.template.metadata.annotations).toBeUndefined();

      const on = parseAllDocuments(
        manifestsFor({ enablePrometheusAnnotations: true, containerPort: 9090 }).deployment,
      )[0].toJS() as Yaml;
      expect(on.spec.template.metadata.annotations).toEqual({
        "prometheus.io/scrape": "true",
        "prometheus.io/port": "9090",
        "prometheus.io/path": METRICS_PATH,
      });
    });

    it("keeps numeric- and boolean-looking names as strings", () => {
      const [deployment] = parseAllDocuments(
        manifestsFor({ workloadName: "123", namespace: "true" }).deployment,
      ).map((d) => d.toJS() as Yaml);
      expect(deployment.metadata.name).toBe("123");
      expect(deployment.metadata.namespace).toBe("true");
    });
  });

  describe("helm values", () => {
    it("mirrors the configuration", () => {
      const values = parse(
        manifestsFor({
          replicas: 2,
          imageTag: "v2",
          containerPort: 9000,
          primaryAddr: "x.y:1234",
          enableProbes: false,
          enablePrometheusAnnotations: true,
        }).helmValues,
      ) as Yaml;

      expect(values.replicaCount).toBe(2);
      expect(values.image.tag).toBe("v2");
      expect(values.service.port).toBe(9000);
      expect(values.faultplane.primaryAddr).toBe("x.y:1234");
      expect(values.faultplane.fallbackAddr).toBeUndefined();
      expect(values.probes.enabled).toBe(false);
      expect(values.prometheus.scrape).toBe(true);
    });
  });
});
