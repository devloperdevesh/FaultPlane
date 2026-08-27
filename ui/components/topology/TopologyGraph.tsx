"use client";

import { useEffect, useState } from "react";

import Connection from "./Connection";
import GatewayNode from "./GatewayNode";
import WorkerNode from "./WorkerNode";

import { getTopology } from "@/lib/api";
import type { TopologySnapshot } from "@/lib/api";

export default function TopologyGraph() {
  const [topology, setTopology] = useState<TopologySnapshot | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    async function loadTopology() {
      try {
        const data = await getTopology();

        if (!mounted) {
          return;
        }

        setTopology(data);
        setError(null);
      } catch (err) {
        if (!mounted) {
          return;
        }

        setError(
          err instanceof Error
            ? err.message
            : "Failed to load topology",
        );
      }
    }

    void loadTopology();

    const interval = window.setInterval(loadTopology, 5000);

    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  if (error) {
    return (
      <div className="rounded-xl border border-red-500/30 bg-zinc-950 p-8">
        <h2 className="text-sm font-semibold text-white">
          Infrastructure Topology
        </h2>

        <p className="mt-4 text-sm text-red-400">
          {error}
        </p>
      </div>
    );
  }

  if (!topology) {
    return (
      <div className="rounded-xl border border-zinc-800 bg-zinc-950 p-8">
        <h2 className="text-sm font-semibold text-white">
          Infrastructure Topology
        </h2>

        <p className="mt-4 text-sm text-zinc-500">
          Loading topology...
        </p>
      </div>
    );
  }

  const gateway = topology.nodes.find(
    (node) => node.type === "gateway",
  );

  const workers = topology.nodes.filter(
    (node) => node.type === "worker",
  );

  const checkpoint = topology.nodes.find(
    (node) => node.type === "checkpoint",
  );

  const gatewayWorkerConnections = topology.connections.filter(
    (connection) =>
      connection.source === gateway?.id ||
      connection.target === gateway?.id,
  );

  const checkpointConnections = topology.connections.filter(
    (connection) =>
      connection.source === checkpoint?.id ||
      connection.target === checkpoint?.id,
  );

  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-950 p-8">
      <div className="mb-6">
        <h2 className="text-sm font-semibold text-white">
          Infrastructure Topology
        </h2>

        <p className="text-xs text-zinc-500">
          Live runtime dependency graph
        </p>
      </div>

      <div className="flex flex-col items-center">
        {gateway && (
          <GatewayNode
            name={gateway.name}
            status={gateway.status}
          />
        )}

        {gatewayWorkerConnections.map((connection) => (
          <Connection
            key={connection.id}
            label={connection.type}
          />
        ))}

        <div className="flex flex-wrap justify-center gap-4">
          {workers.map((worker) => (
            <WorkerNode
              key={worker.id}
              id={worker.name}
              status={worker.status}
            />
          ))}
        </div>

        {checkpointConnections.length > 0 && (
          <>
            <Connection
              label={checkpointConnections[0].type}
            />

            {checkpoint && (
              <div className="rounded-xl border border-emerald-500/30 bg-zinc-950 px-6 py-4 text-center">
                <p className="text-sm font-semibold text-white">
                  {checkpoint.name}
                </p>

                <p className="mt-2 text-xs text-emerald-400">
                  {checkpoint.status}
                </p>
              </div>
            )}
          </>
        )}
      </div>

      <div className="mt-6 border-t border-zinc-800 pt-4">
        <div className="flex items-center justify-between text-xs">
          <span className="text-zinc-500">
            Updated
          </span>

          <span className="text-zinc-400">
            {new Date(topology.updated_at).toLocaleString()}
          </span>
        </div>
      </div>
    </div>
  );
}
