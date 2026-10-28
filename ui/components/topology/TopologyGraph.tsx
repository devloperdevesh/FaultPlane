"use client";

import type { TopologySnapshot } from "@/lib/types";
import { buildImpact, health } from "@/lib/topology/model";
import styles from "./topology.module.css";

type Impact = ReturnType<typeof buildImpact>;

export default function TopologyGraph({
  topology,
  impact,
  selected,
  onSelect,
}: {
  topology: TopologySnapshot;
  impact: Impact;
  selected: string | null;
  onSelect: (id: string) => void;
}) {
  const groups = [...impact.regions.entries()].sort(([a], [b]) =>
    a.localeCompare(b),
  );
  const positions = new Map<string, { x: number; y: number }>();
  const rank = (type: string) =>
    type === "gateway" ? 0 : type === "worker" ? 1 : 2;
  groups.forEach(([, nodes], col) => {
    [...nodes]
      .sort((a, b) => rank(a.type) - rank(b.type) || a.id.localeCompare(b.id))
      .forEach((node, row) =>
        positions.set(node.id, { x: col * 280 + 20, y: row * 110 + 80 }),
      );
  });
  const width = Math.max(560, groups.length * 280);
  const height = Math.max(
    220,
    ...groups.map(([, nodes]) => nodes.length * 110 + 90),
  );

  return (
    <section
      className="rounded-xl border border-zinc-800 bg-zinc-950 p-4"
      aria-labelledby="topology-graph-title"
    >
      <h2 id="topology-graph-title" className="font-semibold text-white">
        Regional topology
      </h2>
      <p className="mt-1 text-sm text-zinc-400">
        Select a node to inspect it. Arrows follow reported source → target
        connections.
      </p>
      <p className="mt-2 text-xs text-zinc-400">
        Red: observed failure · Amber: reachable dependency, not a confirmed
        failure · Green: reported healthy
      </p>
      <div
        className="mt-4 overflow-x-auto rounded-lg border border-zinc-800"
        tabIndex={0}
        aria-label="Scrollable regional topology graph"
      >
        <div className="relative" style={{ width, height }}>
          <svg
            aria-hidden="true"
            width={width}
            height={height}
            className="pointer-events-none absolute inset-0"
          >
            <defs>
              <marker
                id="topology-arrow"
                viewBox="0 0 10 10"
                refX="9"
                refY="5"
                markerWidth="6"
                markerHeight="6"
                orient="auto-start-reverse"
              >
                <path d="M 0 0 L 10 5 L 0 10 z" fill="context-stroke" />
              </marker>
            </defs>
            {topology.connections.map((edge) => {
              const source = positions.get(edge.source);
              const target = positions.get(edge.target);
              if (!source || !target) return null;
              const exposed =
                impact.failed.has(edge.source) ||
                impact.atRisk.has(edge.source);
              const sameColumn = source.x === target.x;
              const rightward = target.x > source.x;
              const startX =
                source.x + (sameColumn ? 120 : rightward ? 240 : 0);
              const startY = source.y + (sameColumn ? 76 : 38);
              const endX = target.x + (sameColumn ? 120 : rightward ? 0 : 240);
              const endY = target.y + (sameColumn ? 0 : 38);
              const path = sameColumn
                ? `M ${startX} ${startY} C ${startX} ${startY + 18}, ${endX} ${endY - 18}, ${endX} ${endY}`
                : `M ${startX} ${startY} C ${startX} ${startY - 80}, ${endX} ${endY - 80}, ${endX} ${endY}`;
              return (
                <path
                  key={`${edge.id}-${impact.latest.get(edge.source)?.timestamp ?? "snapshot"}`}
                  d={path}
                  fill="none"
                  stroke={exposed ? "#fbbf24" : "#64748b"}
                  strokeWidth="2"
                  markerEnd="url(#topology-arrow)"
                  className={exposed ? styles.cascade : undefined}
                />
              );
            })}
          </svg>
          {groups.map(([region, nodes], col) => (
            <div
              key={region}
              className="absolute top-3 w-60 text-center text-xs font-semibold text-zinc-300"
              style={{ left: col * 280 + 20 }}
            >
              {region}
              <span className="mt-1 block font-normal text-zinc-500">
                {nodes.length} nodes ·{" "}
                {nodes.filter((node) => impact.failed.has(node.id)).length}{" "}
                affected
              </span>
            </div>
          ))}
          {topology.nodes.map((node) => {
            const position = positions.get(node.id)!;
            const failed = impact.failed.has(node.id);
            const exposed = impact.atRisk.has(node.id);
            const status = failed
              ? "Observed failure"
              : exposed
                ? "Potential exposure"
                : health(node.status);
            return (
              <button
                key={node.id}
                type="button"
                aria-pressed={selected === node.id}
                aria-controls="node-inspector"
                onClick={() => onSelect(node.id)}
                className={`absolute h-[76px] w-60 rounded-lg border bg-zinc-900 px-3 py-2 text-left focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-sky-400 ${failed ? "border-red-400 text-red-200" : exposed ? "border-amber-400 text-amber-200" : health(node.status) === "healthy" ? "border-emerald-500/60 text-emerald-200" : "border-zinc-700 text-zinc-200"} ${selected === node.id ? "ring-2 ring-sky-400" : ""}`}
                style={{ left: position.x, top: position.y }}
              >
                <span className="block truncate text-sm font-semibold">
                  {node.name}
                </span>
                <span className="block truncate text-xs">
                  {node.type} · {status}
                </span>
                <span className="block truncate text-[11px] text-zinc-400">
                  Reported: {node.status}
                </span>
              </button>
            );
          })}
        </div>
      </div>
    </section>
  );
}
