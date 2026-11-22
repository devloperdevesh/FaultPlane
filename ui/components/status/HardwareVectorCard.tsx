"use client";

import { useEffect, useState } from "react";
import { Cpu, ShieldCheck, Zap } from "lucide-react";
import { getHardwareStatus } from "@/lib/api";
import type { HardwareProfileResponse } from "@/lib/types";

export default function HardwareVectorCard() {
  const [hardware, setHardware] = useState<HardwareProfileResponse | null>(null);

  useEffect(() => {
    let mounted = true;
    const fetchHw = async () => {
      try {
        const data = await getHardwareStatus();
        if (mounted) setHardware(data);
      } catch {
        if (mounted) setHardware(null);
      }
    };
    void fetchHw();
    const interval = window.setInterval(fetchHw, 3000);
    return () => {
      mounted = false;
      window.clearInterval(interval);
    };
  }, []);

  const isPlanetScale = hardware?.active_track?.includes("Planet-Scale") ?? false;

  return (
    <div className="rounded-xl border border-white/10 bg-zinc-950 p-5">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Zap className="h-4 w-4 text-emerald-400" />
          <h3 className="text-sm font-semibold text-white">
            Architecture Engine: Hardware Vector & Fallback Matrix
          </h3>
        </div>
        <span
          className={`rounded-full px-2.5 py-0.5 text-xs font-medium border ${
            isPlanetScale
              ? "border-emerald-500/40 bg-emerald-500/10 text-emerald-300"
              : "border-sky-500/40 bg-sky-500/10 text-sky-300"
          }`}
        >
          {hardware ? hardware.active_track : "Detecting Silicon Profile..."}
        </span>
      </div>

      <p className="mt-2 text-xs text-zinc-400">
        Multi-tiered conditional Layer 4 transport boundary dynamically steered at boot.
      </p>

      <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
        <div className="rounded-lg bg-zinc-900/70 p-3">
          <span className="text-[11px] uppercase tracking-wider text-zinc-500">Vector Registers</span>
          <p className="mt-1 text-base font-semibold text-white">
            {hardware ? `${hardware.vector_bits}-bit` : "64-bit"}
          </p>
          <span className="text-[10px] text-zinc-500">
            {hardware?.vector_bits === 512 ? "Intel ZMM Vector Array" : "Scalar Standard Registers"}
          </span>
        </div>

        <div className="rounded-lg bg-zinc-900/70 p-3">
          <span className="text-[11px] uppercase tracking-wider text-zinc-500">Parallel Lanes</span>
          <p className="mt-1 text-base font-semibold text-white">
            {hardware ? `${hardware.parallel_lanes}x IPs / Cycle` : "1x"}
          </p>
          <span className="text-[10px] text-zinc-500">
            {isPlanetScale ? "16 IP addresses in 1 clock cycle" : "Sequential memory traversal"}
          </span>
        </div>

        <div className="rounded-lg bg-zinc-900/70 p-3">
          <span className="text-[11px] uppercase tracking-wider text-zinc-500">Cgo Boundary</span>
          <p className="mt-1 text-base font-semibold text-white">
            {hardware?.cgo_available ? "Active" : "Isolated"}
          </p>
          <span className="text-[10px] text-zinc-500">
            {hardware?.cgo_available ? "Lock-free unsafe.Pointer" : "Air-gapped safe Go path"}
          </span>
        </div>

        <div className="rounded-lg bg-zinc-900/70 p-3">
          <span className="text-[11px] uppercase tracking-wider text-zinc-500">SLA Bounds</span>
          <p className="mt-1 text-base font-semibold text-emerald-400">&lt; 2.0 ms flat</p>
          <span className="text-[10px] text-zinc-500">Failover hotspot isolation</span>
        </div>
      </div>

      <div className="mt-4 flex flex-wrap items-center gap-2 border-t border-white/5 pt-3">
        <span className="text-xs text-zinc-500">Hardware Flags Evaluated:</span>
        <span
          className={`rounded px-1.5 py-0.5 text-[11px] font-mono ${
            hardware?.has_avx512f ? "bg-emerald-500/20 text-emerald-300" : "bg-zinc-800 text-zinc-400"
          }`}
        >
          AVX-512F: {hardware?.has_avx512f ? "ENABLED" : "ABSENT"}
        </span>
        <span
          className={`rounded px-1.5 py-0.5 text-[11px] font-mono ${
            hardware?.has_avx512vl ? "bg-emerald-500/20 text-emerald-300" : "bg-zinc-800 text-zinc-400"
          }`}
        >
          AVX-512VL: {hardware?.has_avx512vl ? "ENABLED" : "ABSENT"}
        </span>
        <span
          className={`rounded px-1.5 py-0.5 text-[11px] font-mono ${
            hardware?.has_avx512bw ? "bg-emerald-500/20 text-emerald-300" : "bg-zinc-800 text-zinc-400"
          }`}
        >
          AVX-512BW: {hardware?.has_avx512bw ? "ENABLED" : "ABSENT"}
        </span>
      </div>
    </div>
  );
}
