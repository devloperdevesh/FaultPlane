"use client";

import { useMemo, useRef, useState } from "react";
import type { ChangeEvent, ReactNode } from "react";

import Panel from "@/components/ui/Panel";
import {
  DEFAULT_CONFIG,
  LOG_LEVELS,
  type InstallerConfig,
  type LogLevel,
} from "@/lib/installer/config";
import {
  FILE_NAMES,
  generateManifests,
  type GeneratedManifests,
} from "@/lib/installer/generate";
import type { FieldName } from "@/lib/installer/validate";

type OutputKey = keyof GeneratedManifests;

const OUTPUT_TABS: { key: OutputKey; label: string }[] = [
  { key: "sidecar", label: "Sidecar container" },
  { key: "deployment", label: "Deployment + Service" },
  { key: "helmValues", label: "Helm values" },
];

const inputClass =
  "w-full rounded-lg border border-zinc-800 bg-black/40 px-3 py-2 text-sm text-white placeholder:text-zinc-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 aria-[invalid=true]:border-red-500/70";

interface FieldProps {
  id: FieldName;
  label: string;
  hint?: string;
  error?: string;
  children: ReactNode;
}

function Field({ id, label, hint, error, children }: FieldProps) {
  return (
    <div className="space-y-1.5">
      <label htmlFor={id} className="text-xs font-medium text-zinc-300">
        {label}
      </label>
      {children}
      {error ? (
        <p id={`${id}-error`} role="alert" className="text-xs text-red-400">
          {error}
        </p>
      ) : hint ? (
        <p className="text-xs text-zinc-500">{hint}</p>
      ) : null}
    </div>
  );
}

function Toggle({
  id,
  label,
  checked,
  onChange,
}: {
  id: string;
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
}) {
  return (
    <label htmlFor={id} className="flex items-center gap-3 text-sm text-zinc-300">
      <input
        id={id}
        type="checkbox"
        checked={checked}
        onChange={(event) => onChange(event.target.checked)}
        className="h-4 w-4 rounded border-zinc-700 bg-black accent-emerald-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400"
      />
      {label}
    </label>
  );
}

export default function ManifestInstaller() {
  const [config, setConfig] = useState<InstallerConfig>(DEFAULT_CONFIG);
  const [activeTab, setActiveTab] = useState<OutputKey>("sidecar");
  const [status, setStatus] = useState("");
  const statusTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const result = useMemo(() => generateManifests(config), [config]);

  const update = <K extends keyof InstallerConfig>(key: K, value: InstallerConfig[K]) => {
    setConfig((current) => ({ ...current, [key]: value }));
  };

  const text =
    (key: keyof InstallerConfig) => (event: ChangeEvent<HTMLInputElement>) =>
      update(key, event.target.value as never);

  const numeric =
    (key: "replicas" | "containerPort") => (event: ChangeEvent<HTMLInputElement>) =>
      update(key, event.target.value === "" ? Number.NaN : Number(event.target.value));

  const errors = result.ok ? {} : result.errors;
  const errorCount = Object.keys(errors).length;

  const fieldProps = (id: FieldName) => ({
    id,
    "aria-invalid": errors[id] ? (true as const) : undefined,
    "aria-describedby": errors[id] ? `${id}-error` : undefined,
    className: inputClass,
  });

  const flash = (message: string) => {
    setStatus(message);
    if (statusTimer.current) clearTimeout(statusTimer.current);
    statusTimer.current = setTimeout(() => setStatus(""), 2500);
  };

  const activeText = result.ok ? result.manifests[activeTab] : "";

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(activeText);
      flash("Copied to clipboard.");
    } catch {
      flash("Copy failed. Select the text and copy it manually.");
    }
  };

  const download = () => {
    const url = URL.createObjectURL(new Blob([activeText], { type: "text/yaml" }));
    const link = document.createElement("a");
    link.href = url;
    link.download = FILE_NAMES[activeTab];
    document.body.appendChild(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
    flash(`Downloaded ${FILE_NAMES[activeTab]}.`);
  };

  return (
    <div className="space-y-6">
      <div
        role="note"
        className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-4 text-sm text-amber-200"
      >
        <p className="font-semibold">Preview: manifest generator only</p>
        <p className="mt-1 text-amber-200/80">
          This page generates files from the values you enter. It does not connect to a cluster, deploy
          anything or read live cluster data. Review every manifest before applying it.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-6 xl:grid-cols-2">
        <div className="min-w-0 space-y-6">
          <Panel>
            <h2 className="text-sm font-semibold text-white">Workload</h2>
            <div className="mt-5 grid gap-4 sm:grid-cols-2">
              <Field id="namespace" label="Namespace" error={errors.namespace} hint="Must already exist.">
                <input {...fieldProps("namespace")} value={config.namespace} onChange={text("namespace")} />
              </Field>
              <Field id="workloadName" label="Workload name" error={errors.workloadName}>
                <input {...fieldProps("workloadName")} value={config.workloadName} onChange={text("workloadName")} />
              </Field>
              <Field id="replicas" label="Replicas" error={errors.replicas} hint="Standalone Deployment only.">
                <input
                  {...fieldProps("replicas")}
                  type="number"
                  inputMode="numeric"
                  value={Number.isNaN(config.replicas) ? "" : config.replicas}
                  onChange={numeric("replicas")}
                />
              </Field>
              <Field id="containerPort" label="Container port" error={errors.containerPort}>
                <input
                  {...fieldProps("containerPort")}
                  type="number"
                  inputMode="numeric"
                  value={Number.isNaN(config.containerPort) ? "" : config.containerPort}
                  onChange={numeric("containerPort")}
                />
              </Field>
            </div>
          </Panel>

          <Panel>
            <h2 className="text-sm font-semibold text-white">Image</h2>
            <div className="mt-5 grid gap-4 sm:grid-cols-2">
              <Field
                id="imageRepository"
                label="Image repository"
                error={errors.imageRepository}
                hint="Editable placeholder. Use the image you publish."
              >
                <input {...fieldProps("imageRepository")} value={config.imageRepository} onChange={text("imageRepository")} />
              </Field>
              <Field id="imageTag" label="Image tag" error={errors.imageTag}>
                <input {...fieldProps("imageTag")} value={config.imageTag} onChange={text("imageTag")} />
              </Field>
              <Field id="logLevel" label="Log level" error={errors.logLevel}>
                <select
                  id="logLevel"
                  className={inputClass}
                  value={config.logLevel}
                  onChange={(event) => update("logLevel", event.target.value as LogLevel)}
                >
                  {LOG_LEVELS.map((level) => (
                    <option key={level} value={level}>
                      {level}
                    </option>
                  ))}
                </select>
              </Field>
            </div>
          </Panel>

          <Panel>
            <h2 className="text-sm font-semibold text-white">Upstreams (optional)</h2>
            <p className="mt-2 text-xs text-zinc-500">
              Written as FAULTPLANE_PRIMARY_ADDR and FAULTPLANE_FALLBACK_ADDR, as in the repository compose example.
              Leave empty to omit them.
            </p>
            <div className="mt-4 grid gap-4 sm:grid-cols-2">
              <Field id="primaryAddr" label="Primary (host:port)" error={errors.primaryAddr}>
                <input
                  {...fieldProps("primaryAddr")}
                  placeholder="backend.default.svc:8080"
                  value={config.primaryAddr}
                  onChange={text("primaryAddr")}
                />
              </Field>
              <Field id="fallbackAddr" label="Fallback (host:port)" error={errors.fallbackAddr}>
                <input
                  {...fieldProps("fallbackAddr")}
                  placeholder="backup.default.svc:8080"
                  value={config.fallbackAddr}
                  onChange={text("fallbackAddr")}
                />
              </Field>
            </div>
          </Panel>

          <Panel>
            <h2 className="text-sm font-semibold text-white">Resources</h2>
            <div className="mt-5 grid gap-4 sm:grid-cols-2">
              <Field id="cpuRequest" label="CPU request" error={errors.cpuRequest}>
                <input {...fieldProps("cpuRequest")} value={config.cpuRequest} onChange={text("cpuRequest")} />
              </Field>
              <Field id="cpuLimit" label="CPU limit" error={errors.cpuLimit}>
                <input {...fieldProps("cpuLimit")} value={config.cpuLimit} onChange={text("cpuLimit")} />
              </Field>
              <Field id="memoryRequest" label="Memory request" error={errors.memoryRequest}>
                <input {...fieldProps("memoryRequest")} value={config.memoryRequest} onChange={text("memoryRequest")} />
              </Field>
              <Field id="memoryLimit" label="Memory limit" error={errors.memoryLimit}>
                <input {...fieldProps("memoryLimit")} value={config.memoryLimit} onChange={text("memoryLimit")} />
              </Field>
            </div>
          </Panel>

          <Panel>
            <h2 className="text-sm font-semibold text-white">Options</h2>
            <div className="mt-4 space-y-3">
              <Toggle
                id="enableProbes"
                label="Readiness and liveness probes on /health"
                checked={config.enableProbes}
                onChange={(value) => update("enableProbes", value)}
              />
              <Toggle
                id="enablePrometheusAnnotations"
                label="Prometheus scrape annotations"
                checked={config.enablePrometheusAnnotations}
                onChange={(value) => update("enablePrometheusAnnotations", value)}
              />
            </div>
            <p className="mt-4 text-xs text-zinc-500">
              Generated containers are unprivileged: no privileged mode and no added capabilities. Privileged eBPF
              settings are not generated by this installer.
            </p>
            <button
              type="button"
              onClick={() => setConfig(DEFAULT_CONFIG)}
              className="mt-5 rounded-lg border border-zinc-700 px-3 py-1.5 text-xs font-medium text-zinc-300 transition-colors hover:border-zinc-500 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400"
            >
              Reset to defaults
            </button>
          </Panel>
        </div>

        <div className="min-w-0 space-y-6">
          <Panel>
            <div className="flex flex-wrap items-center justify-between gap-3">
              <h2 className="text-sm font-semibold text-white">Generated output</h2>
              <p className="text-xs text-zinc-500">
                Cluster connection: <span className="text-zinc-300">not connected (preview)</span>
              </p>
            </div>

            <div role="tablist" aria-label="Generated output" className="mt-4 flex flex-wrap gap-2">
              {OUTPUT_TABS.map((tab) => (
                <button
                  key={tab.key}
                  id={`tab-${tab.key}`}
                  role="tab"
                  type="button"
                  aria-selected={activeTab === tab.key}
                  aria-controls="output-panel"
                  onClick={() => setActiveTab(tab.key)}
                  className={`rounded-full px-3 py-1.5 text-xs font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 ${
                    activeTab === tab.key
                      ? "bg-emerald-500/20 text-emerald-300"
                      : "text-zinc-400 hover:bg-white/5 hover:text-white"
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </div>

            <div id="output-panel" role="tabpanel" aria-labelledby={`tab-${activeTab}`} className="mt-4">
              {result.ok ? (
                <>
                  {result.warnings.length > 0 && (
                    <ul className="mb-3 space-y-1 text-xs text-amber-300">
                      {result.warnings.map((warning) => (
                        <li key={warning}>Warning: {warning}</li>
                      ))}
                    </ul>
                  )}
                  <pre
                    tabIndex={0}
                    aria-label={`${FILE_NAMES[activeTab]} preview`}
                    className="max-h-[520px] overflow-auto rounded-lg border border-zinc-800 bg-black/50 p-4 text-xs leading-5 text-zinc-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400"
                  >
                    <code>{activeText}</code>
                  </pre>
                  <div className="mt-4 flex flex-wrap items-center gap-3">
                    <button
                      type="button"
                      onClick={copy}
                      className="rounded-lg bg-emerald-500 px-4 py-2 text-sm font-semibold text-black transition-colors hover:bg-emerald-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-300"
                    >
                      Copy
                    </button>
                    <button
                      type="button"
                      onClick={download}
                      className="rounded-lg border border-zinc-700 px-4 py-2 text-sm font-medium text-zinc-200 transition-colors hover:border-zinc-500 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400"
                    >
                      Download {FILE_NAMES[activeTab]}
                    </button>
                    <span role="status" aria-live="polite" className="text-xs text-zinc-400">
                      {status}
                    </span>
                  </div>
                </>
              ) : (
                <div
                  role="alert"
                  className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-300"
                >
                  Fix {errorCount} {errorCount === 1 ? "field" : "fields"} to generate the manifest. The errors are
                  shown next to each input.
                </div>
              )}
            </div>
          </Panel>

          {result.ok && (
            <Panel>
              <h2 className="text-sm font-semibold text-white">Validate and apply</h2>
              <p className="mt-2 text-xs text-zinc-500">
                Run these yourself. This page never runs them for you.
              </p>
              <pre className="mt-3 overflow-x-auto rounded-lg border border-zinc-800 bg-black/50 p-4 text-xs leading-5 text-zinc-200">
                <code>
                  {`kubectl create namespace ${config.namespace}   # if it does not exist yet\n` +
                    `kubectl apply --dry-run=client -f ${FILE_NAMES.deployment}\n` +
                    `kubectl apply -f ${FILE_NAMES.deployment}`}
                </code>
              </pre>
            </Panel>
          )}
        </div>
      </div>
    </div>
  );
}
