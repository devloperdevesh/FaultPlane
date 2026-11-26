"use client";

import { useCallback, useEffect, useState } from "react";
import {
  Cpu,
  Lock,
  Activity,
  ShieldCheck,
  RefreshCw,
  Zap,
  CheckCircle2,
  AlertTriangle,
  Database,
} from "lucide-react";
import MemoryMatrix3D from "./MemoryMatrix3D";
import { getMemoryMatrix } from "@/lib/api";
import type { BufferPoolSnapshot, SlotSnapshot } from "@/lib/types";

function buildFallbackSnapshot(): BufferPoolSnapshot {
  const mockSlots: SlotSnapshot[] = Array.from({ length: 64 }).map((_, i) => {
    const domain = (i % 4) + 1;
    const isFenced = i % 8 === 0;
    const isActive = i % 3 !== 0;

    let status = "FREE";
    let color = "#1E293B";

    if (isFenced) {
      status = "FENCED";
      color = "#FF0055";
    } else if (isActive) {
      status = "ACTIVE";
      switch (domain) {
        case 1:
          color = "#00F0FF";
          break;
        case 2:
          color = "#00FF88";
          break;
        case 3:
          color = "#BD00FF";
          break;
        default:
          color = "#FFB800";
      }
    }

    return {
      index: i,
      slot_id: i,
      tenant_id: status === "FREE" ? "000000000000" : `a7f9b2c${i.toString(16).padStart(5, "0")}`,
      domain_key: status === "FREE" ? 0 : domain,
      status,
      active_tokens: isActive ? 24 + ((i * 3) % 40) : 0,
      access_count: status === "FREE" ? 0 : (i + 1) * 1420,
      memory_offset_hex: `0x${(i * 64).toString(16).toUpperCase().padStart(4, "0")}`,
      cache_line_aligned: true,
      color_channel: color,
    };
  });

  const activeSlots = mockSlots.filter((s) => s.status === "ACTIVE").length;
  const fencedSlots = mockSlots.filter((s) => s.status === "FENCED").length;

  return {
    total_slots: 64,
    active_slots: activeSlots,
    fenced_slots: fencedSlots,
    cache_line_size: 64,
    cache_alignment_congruent: true,
    l1_l2_thrashing_rate: null,
    l1_l2_thrashing_available: false,
    l1_l2_thrashing_note: "Backend stream offline — hardware PMU counters unavailable in userspace.",
    gc_pauses_ms: 0,
    gc_pause_count: 0,
    gc_pauses_available: false,
    isolation_mode: "logical_domain_tag",
    hardware_mpk_enforced: false,
    active_tenants_count: 4,
    slots: mockSlots,
  };
}

function pickSelectedSlot(slots: SlotSnapshot[], prev: SlotSnapshot | null): SlotSnapshot | null {
  if (prev) {
    return slots.find((s) => s.slot_id === prev.slot_id) ?? slots[0] ?? null;
  }
  return slots[0] ?? null;
}

export default function MemoryCockpitOverview() {
  const [snapshot, setSnapshot] = useState<BufferPoolSnapshot | null>(null);
  const [loading, setLoading] = useState(true);
  const [offlineFallback, setOfflineFallback] = useState(false);
  const [selectedSlot, setSelectedSlot] = useState<SlotSnapshot | null>(null);

  const applyLive = useCallback((data: BufferPoolSnapshot) => {
    setSnapshot(data);
    setOfflineFallback(false);
    setSelectedSlot((prev) => pickSelectedSlot(data.slots, prev));
  }, []);

  const applyOffline = useCallback(() => {
    const mockSnap = buildFallbackSnapshot();
    setSnapshot(mockSnap);
    setOfflineFallback(true);
    setSelectedSlot((prev) => pickSelectedSlot(mockSnap.slots, prev));
  }, []);

  const fetchMatrix = async () => {
    try {
      const data = await getMemoryMatrix();
      applyLive(data);
    } catch {
      applyOffline();
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let cancelled = false;
    const poll = () => {
      getMemoryMatrix()
        .then((data) => {
          if (!cancelled) applyLive(data);
        })
        .catch(() => {
          if (!cancelled) applyOffline();
        })
        .finally(() => {
          if (!cancelled) setLoading(false);
        });
    };

    poll();
    const interval = window.setInterval(poll, 2000);
    return () => {
      cancelled = true;
      window.clearInterval(interval);
    };
  }, [applyLive, applyOffline]);

  const alignedCount = snapshot?.slots.filter((s) => s.cache_line_aligned).length ?? 0;

  return (
    <div className="space-y-6">
      {/* 1. Cockpit Header Banner */}
      <div className="rounded-2xl border border-white/10 bg-gradient-to-r from-zinc-950 via-zinc-900 to-zinc-950 p-6 shadow-xl backdrop-blur-xl">
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="space-y-1">
            <div className="flex flex-wrap items-center gap-2">
              <span className="flex h-6 items-center gap-1.5 rounded-full border border-cyan-500/30 bg-cyan-500/10 px-2.5 text-xs font-medium text-cyan-400">
                <Cpu className="h-3.5 w-3.5" />
                Cache-Line Aligned Buffer Pool
              </span>
              <span
                className={`flex h-6 items-center gap-1 rounded-full border px-2 text-xs font-medium ${
                  snapshot?.cache_alignment_congruent
                    ? "border-emerald-500/30 bg-emerald-500/10 text-emerald-400"
                    : "border-amber-500/30 bg-amber-500/10 text-amber-400"
                }`}
              >
                {snapshot?.cache_alignment_congruent ? (
                  <CheckCircle2 className="h-3 w-3" />
                ) : (
                  <AlertTriangle className="h-3 w-3" />
                )}
                {snapshot
                  ? snapshot.cache_alignment_congruent
                    ? "64-Byte Address Congruent"
                    : "Alignment Unverified"
                  : "Checking Alignment"}
              </span>
              {offlineFallback && (
                <span className="flex h-6 items-center gap-1 rounded-full border border-amber-500/30 bg-amber-500/10 px-2.5 text-xs font-medium text-amber-300">
                  <AlertTriangle className="h-3 w-3" />
                  Backend Offline (Preview Fallback)
                </span>
              )}
            </div>
            <h2 className="text-xl font-bold tracking-tight text-white">
              Multi-Tenant Memory Cockpit &amp; Structural Matrix
            </h2>
            <p className="text-sm text-zinc-400">
              64-byte aligned slot ring with lock-free CAS state transitions, measured Go runtime GC telemetry, and logical domain fencing (software isolation mode; hardware Intel MPK/PKRU not active).
            </p>
          </div>

          <button
            onClick={fetchMatrix}
            className="inline-flex items-center gap-1.5 rounded-lg border border-white/10 bg-zinc-900 px-3.5 py-2 text-xs font-medium text-zinc-300 transition hover:bg-zinc-800"
          >
            <RefreshCw className={`h-3.5 w-3.5 ${loading ? "animate-spin" : ""}`} />
            Refresh Matrix Snapshot
          </button>
        </div>
      </div>

      {/* 2. Buffer Pool Telemetry KPI Cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {/* KPI 1: Structural Cache-Line Alignment */}
        <div className="rounded-xl border border-white/10 bg-zinc-950/80 p-4 backdrop-blur-md">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-zinc-400">Cache Line Alignment</span>
            <ShieldCheck
              className={`h-4 w-4 ${
                snapshot?.cache_alignment_congruent ? "text-emerald-400" : "text-amber-400"
              }`}
            />
          </div>
          <p className="mt-2 text-2xl font-bold tracking-tight text-white">
            {snapshot
              ? snapshot.cache_alignment_congruent
                ? "100% Congruent"
                : "Misaligned"
              : "Loading..."}
          </p>
          <p className="mt-1 text-xs text-emerald-400/90 font-medium">
            {alignedCount}/{snapshot?.total_slots ?? 0} slots aligned on {snapshot?.cache_line_size ?? 64}-byte boundaries
          </p>
        </div>

        {/* KPI 2: L1/L2 Cache Thrashing (requires PMU telemetry) */}
        <div className="rounded-xl border border-white/10 bg-zinc-950/80 p-4 backdrop-blur-md">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-zinc-400">L1/L2 Cache Thrashing</span>
            <Activity className="h-4 w-4 text-cyan-400" />
          </div>
          <p className="mt-2 text-2xl font-bold tracking-tight text-white">
            {snapshot?.l1_l2_thrashing_available && snapshot.l1_l2_thrashing_rate !== null
              ? `${snapshot.l1_l2_thrashing_rate.toFixed(2)}%`
              : "Unavailable"}
          </p>
          <p className="mt-1 text-xs text-zinc-400">
            {snapshot?.l1_l2_thrashing_note ?? "Requires hardware PMU counters (unavailable in userspace)"}
          </p>
        </div>

        {/* KPI 3: Go Runtime GC Pauses */}
        <div className="rounded-xl border border-white/10 bg-zinc-950/80 p-4 backdrop-blur-md">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-zinc-400">Runtime GC Pauses (Total)</span>
            <Zap className="h-4 w-4 text-emerald-400" />
          </div>
          <p className="mt-2 text-2xl font-bold tracking-tight text-white">
            {snapshot?.gc_pauses_available
              ? `${snapshot.gc_pauses_ms.toFixed(3)} ms`
              : "Unavailable"}
          </p>
          <p className="mt-1 text-xs text-zinc-400">
            {snapshot?.gc_pauses_available
              ? `Measured Go STW GC pauses across ${snapshot.gc_pause_count} cycles`
              : "Live runtime GC telemetry unavailable offline"}
          </p>
        </div>

        {/* KPI 4: Logical Domain Fencing */}
        <div className="rounded-xl border border-white/10 bg-zinc-950/80 p-4 backdrop-blur-md">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium text-zinc-400">Logical Domain Fences</span>
            <Lock className="h-4 w-4 text-rose-400" />
          </div>
          <p className="mt-2 text-2xl font-bold tracking-tight text-white">
            {snapshot?.fenced_slots ?? 0} Fenced Slots
          </p>
          <p className="mt-1 text-xs text-rose-400/90 font-medium">
            {snapshot?.hardware_mpk_enforced
              ? "Hardware MPK/PKRU Enforced"
              : "Logical software domain tag only (MPK inactive)"}
          </p>
        </div>
      </div>

      {/* 3. Isometric Matrix Visualizer */}
      <MemoryMatrix3D
        snapshot={snapshot}
        selectedSlot={selectedSlot}
        onSelectSlot={setSelectedSlot}
      />

      {/* 4. Multi-Tenant Buffer Pool Layout Table */}
      <div className="overflow-hidden rounded-xl border border-white/10 bg-zinc-950">
        <div className="border-b border-white/10 px-5 py-3.5 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Database className="h-4 w-4 text-cyan-400" />
            <h3 className="text-sm font-semibold text-white">
              Multi-Tenant Buffer Pool Layout ({snapshot?.slots.length ?? 0} Slots)
            </h3>
          </div>
          <span className="text-xs text-zinc-400">
            Click any row to inspect slot register
          </span>
        </div>

        <div className="overflow-x-auto max-h-72">
          <table className="w-full text-left text-xs">
            <thead className="sticky top-0 border-b border-white/10 bg-zinc-900/95 backdrop-blur text-zinc-400 font-medium">
              <tr>
                <th className="px-4 py-2.5">Slot ID</th>
                <th className="px-4 py-2.5">Memory Offset</th>
                <th className="px-4 py-2.5">64B Alignment</th>
                <th className="px-4 py-2.5">Tenant Workspace Token</th>
                <th className="px-4 py-2.5">Domain Tag</th>
                <th className="px-4 py-2.5">Status</th>
                <th className="px-4 py-2.5">Active Tokens</th>
                <th className="px-4 py-2.5">Acquisitions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {snapshot?.slots.map((slot) => {
                const isSelected = selectedSlot?.slot_id === slot.slot_id;
                return (
                  <tr
                    key={slot.slot_id}
                    onClick={() => setSelectedSlot(slot)}
                    className={`cursor-pointer transition ${
                      isSelected ? "bg-cyan-500/10" : "hover:bg-zinc-900/40"
                    }`}
                  >
                    <td className="px-4 py-2.5 font-mono text-zinc-200">
                      #{slot.slot_id}
                    </td>
                    <td className="px-4 py-2.5 font-mono text-cyan-400">
                      {slot.memory_offset_hex}
                    </td>
                    <td className="px-4 py-2.5">
                      {slot.cache_line_aligned ? (
                        <span className="inline-flex items-center gap-1 text-[11px] font-medium text-emerald-400">
                          <CheckCircle2 className="h-3 w-3" />
                          64-Byte Aligned
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1 text-[11px] font-medium text-rose-400">
                          <AlertTriangle className="h-3 w-3" />
                          Unaligned
                        </span>
                      )}
                    </td>
                    <td className="px-4 py-2.5 font-mono text-zinc-300">
                      {slot.tenant_id}
                    </td>
                    <td className="px-4 py-2.5">
                      <span className="font-mono text-white">
                        Domain 0x0{slot.domain_key}
                      </span>
                    </td>
                    <td className="px-4 py-2.5">
                      <span
                        className="rounded px-2 py-0.5 text-[10px] font-bold"
                        style={{
                          backgroundColor: `${slot.color_channel}20`,
                          color: slot.color_channel,
                          border: `1px solid ${slot.color_channel}50`,
                        }}
                      >
                        {slot.status}
                      </span>
                    </td>
                    <td className="px-4 py-2.5 text-zinc-300">
                      {slot.active_tokens} tokens
                    </td>
                    <td className="px-4 py-2.5 font-mono text-zinc-400">
                      {slot.access_count.toLocaleString()}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
