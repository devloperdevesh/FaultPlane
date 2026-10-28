"use client";

import { useEffect, useState } from "react";
import {
  buildImpact,
  eventKind,
  eventTargets,
  regionOf,
} from "@/lib/topology/model";
import { loadTopologyData, type TopologyData } from "@/lib/topology/load";
import { demoSnapshot } from "@/lib/topology/demo";
import TopologyGraph from "./TopologyGraph";

export function timestamp(value: string) {
  const date = new Date(value);
  return Number.isFinite(date.getTime())
    ? date.toISOString()
    : "Timestamp not reported";
}

export default function TopologyDashboard({
  demo = false,
}: {
  demo?: boolean;
}) {
  const [data, setData] = useState<TopologyData | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [receivedAt, setReceivedAt] = useState<string | null>(null);
  const [selected, setSelected] = useState<string | null>(null);
  const [retry, setRetry] = useState(0);
  const [step, setStep] = useState(1);

  useEffect(() => {
    if (demo) return;
    let disposed = false;
    let timer: ReturnType<typeof setTimeout>;
    let controller: AbortController;
    async function poll() {
      controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 10000);
      try {
        const next = await loadTopologyData(controller.signal);
        if (!disposed) {
          setData(next);
          setError(null);
          setReceivedAt(new Date().toISOString());
        }
      } catch (err) {
        if (!disposed)
          setError(
            err instanceof Error
              ? err.message
              : "Unable to reach the FaultPlane API",
          );
      } finally {
        clearTimeout(timeout);
        controller.abort();
        if (!disposed) timer = setTimeout(poll, 5000);
      }
    }
    void poll();
    return () => {
      disposed = true;
      clearTimeout(timer);
      controller?.abort();
    };
  }, [demo, retry]);

  const current = demo ? demoSnapshot(step) : data;
  const impact = current ? buildImpact(current.topology, current.events) : null;
  const node = current?.topology.nodes.find((item) => item.id === selected);
  const worker = current?.workers.find((item) => item.id === selected);
  const latest = node ? impact?.latest.get(node.id) : undefined;

  return (
    <div className="space-y-6">
      {demo && (
        <section
          className="rounded-xl border border-amber-400/40 bg-amber-400/10 p-4"
          aria-label="Demo controls"
        >
          <p className="font-semibold text-amber-200">
            DEMO — fixture replay, no live telemetry
          </p>
          <p className="mt-1 text-sm text-amber-100/80">
            Two fictional regions. Replay one worker failure and observed
            recovery. No runtime changes are made.
          </p>
          <div className="mt-3 flex flex-wrap gap-2">
            {["1. Healthy", "2. Failure", "3. Recovered"].map(
              (label, index) => (
                <button
                  key={label}
                  type="button"
                  aria-pressed={step === index}
                  onClick={() => setStep(index)}
                  className="rounded border border-amber-300/40 px-3 py-2 text-sm text-amber-100 aria-pressed:bg-amber-300/20 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-200"
                >
                  {label}
                </button>
              ),
            )}
          </div>
        </section>
      )}

      {!demo && (
        <div
          role="status"
          className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-zinc-800 bg-zinc-950 p-4 text-sm"
        >
          <div>
            <p className={error ? "text-amber-300" : "text-emerald-300"}>
              {error
                ? data
                  ? "Connection interrupted — stale snapshot"
                  : "Disconnected / API error"
                : data
                  ? "Connected · polling every 5 seconds"
                  : "Loading runtime telemetry…"}
            </p>
            <p className="mt-1 text-xs text-zinc-400">
              {receivedAt
                ? `Last successful fetch: ${timestamp(receivedAt)}`
                : "Waiting for the first valid API response."}
            </p>
            {error && (
              <p className="mt-2 max-w-3xl break-words text-xs text-amber-200">
                {error}. Retrying automatically. Check the daemon and
                NEXT_PUBLIC_API_URL.
              </p>
            )}
          </div>
          <button
            type="button"
            onClick={() => setRetry((value) => value + 1)}
            className="rounded border border-zinc-600 px-3 py-2 text-zinc-200 focus-visible:outline-2 focus-visible:outline-sky-400"
          >
            Retry now
          </button>
        </div>
      )}

      {current && impact && (
        <>
          <div className="grid grid-cols-2 gap-3 xl:grid-cols-4">
            {[
              [
                "Reported regions",
                [...impact.regions.keys()].filter(
                  (region) => region !== "Region not reported",
                ).length,
              ],
              [
                "Worker nodes",
                current.topology.nodes.filter((item) => item.type === "worker")
                  .length,
              ],
              ["Observed failures", impact.failed.size],
              ["Potentially exposed", impact.atRisk.size],
            ].map(([label, value]) => (
              <div
                key={label}
                className="rounded-xl border border-zinc-800 bg-zinc-950 p-4"
              >
                <p className="text-xs text-zinc-400">{label}</p>
                <p className="mt-1 text-2xl font-semibold text-white">
                  {value}
                </p>
              </div>
            ))}
          </div>
          <p className="text-xs text-zinc-400">
            Snapshot: {timestamp(current.topology.updated_at)}. Potential
            exposure follows directed edges; it does not establish causality or
            confirm additional failures.
          </p>
          {current.topology.nodes.some((item) => !item.region?.trim()) && (
            <p className="rounded-lg border border-zinc-700 p-3 text-sm text-zinc-400">
              Some nodes have no geographic metadata. They appear under “Region
              not reported”; regions are never inferred from names or IPs.
            </p>
          )}
          {current.topology.nodes.length === 0 ? (
            <p className="rounded-xl border border-zinc-800 p-6 text-zinc-400">
              No topology nodes reported by the runtime.
            </p>
          ) : (
            <TopologyGraph
              topology={current.topology}
              impact={impact}
              selected={selected}
              onSelect={setSelected}
            />
          )}

          <div className="grid gap-6 xl:grid-cols-2">
            <section
              id="node-inspector"
              className="scroll-mt-4 rounded-xl border border-zinc-800 bg-zinc-950 p-5"
              aria-labelledby="inspector-title"
            >
              <h2 id="inspector-title" className="font-semibold text-white">
                Node inspector
              </h2>
              {node ? (
                <dl className="mt-4 grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 break-words text-sm">
                  <dt className="text-zinc-400">Node</dt>
                  <dd>
                    {node.name} ({node.id})
                  </dd>
                  <dt className="text-zinc-400">Type / region</dt>
                  <dd>
                    {node.type} / {regionOf(node)}
                  </dd>
                  <dt className="text-zinc-400">Reported status</dt>
                  <dd>{node.status}</dd>
                  <dt className="text-zinc-400">CPU / memory</dt>
                  <dd>
                    {worker
                      ? `${worker.cpu.toFixed(1)}% / ${(worker.memory / 1024 / 1024).toFixed(1)} MiB`
                      : "Not reported for this node"}
                  </dd>
                  <dt className="text-zinc-400">Last incident</dt>
                  <dd>
                    {latest
                      ? `${latest.type} · ${timestamp(latest.timestamp)}`
                      : "No correlated incident in the retained event window"}
                  </dd>
                  <dt className="text-zinc-400">Exposure</dt>
                  <dd>
                    {impact.failed.has(node.id)
                      ? "Observed failure"
                      : impact.atRisk.has(node.id)
                        ? "Reachable from an observed failure; failure unconfirmed"
                        : "No exposure found in reported paths"}
                  </dd>
                </dl>
              ) : (
                <p className="mt-4 text-sm text-zinc-400">
                  {selected
                    ? "Selected node is no longer in this snapshot. Select another node."
                    : "Select a node in the graph to inspect its runtime state."}
                </p>
              )}
              {node && (
                <ul
                  className="mt-4 space-y-2 text-xs text-zinc-400"
                  aria-label="Selected node connections"
                >
                  {current.topology.connections
                    .filter(
                      (edge) =>
                        edge.source === node.id || edge.target === node.id,
                    )
                    .map((edge) => (
                      <li key={edge.id} className="break-words">
                        {edge.source} → {edge.target} · {edge.type} ·{" "}
                        {edge.status}
                      </li>
                    ))}
                </ul>
              )}
            </section>
            <section
              className="rounded-xl border border-zinc-800 bg-zinc-950 p-5"
              aria-labelledby="timeline-title"
            >
              <h2 id="timeline-title" className="font-semibold text-white">
                Failure & recovery timeline
              </h2>
              <p className="mt-1 text-xs text-zinc-400">
                Latest 50 incident records in the API’s retained window. Generic
                HTTP requests are excluded.
              </p>
              {impact.timeline.length === 0 ? (
                <p className="mt-4 text-sm text-zinc-400">
                  No failure or recovery events recorded. Request activity alone
                  does not indicate a healthy or failed connection.
                </p>
              ) : (
                <ol className="mt-4 max-h-96 space-y-3 overflow-y-auto">
                  {impact.timeline
                    .slice(-50)
                    .reverse()
                    .map((event, index) => {
                      const targets = eventTargets(event, current.topology);
                      return (
                        <li
                          key={`${event.type}-${event.timestamp}-${index}`}
                          className={`rounded-lg border p-3 text-sm ${eventKind(event) === "failure" ? "border-red-400/30 text-red-200" : "border-emerald-400/30 text-emerald-200"}`}
                        >
                          <p className="font-medium">
                            {event.type} ·{" "}
                            {eventKind(event) === "recovery"
                              ? "Recovery observed"
                              : "Failure reported"}
                          </p>
                          <time
                            dateTime={event.timestamp}
                            className="mt-1 block text-xs text-zinc-400"
                          >
                            {timestamp(event.timestamp)}
                          </time>
                          <p className="mt-1 break-words text-xs">
                            {targets.length
                              ? `Affected nodes: ${targets.join(", ")}`
                              : "Uncorrelated event — affected nodes and blast radius unknown"}
                          </p>
                          {event.metadata?.source && (
                            <p className="mt-1 break-words text-xs text-zinc-400">
                              Source: {event.metadata.source}
                            </p>
                          )}
                        </li>
                      );
                    })}
                </ol>
              )}
            </section>
          </div>
        </>
      )}
    </div>
  );
}
