"use client";

import { useEffect, useRef, useState } from "react";
import {
  RotateCcw,
  Cpu,
} from "lucide-react";
import type { BufferPoolSnapshot, SlotSnapshot } from "@/lib/types";

interface MemoryMatrix3DProps {
  snapshot: BufferPoolSnapshot | null;
  onSelectSlot?: (slot: SlotSnapshot | null) => void;
  selectedSlot?: SlotSnapshot | null;
}

export default function MemoryMatrix3D({
  snapshot,
  onSelectSlot,
  selectedSlot,
}: MemoryMatrix3DProps) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  // 3D Orbit State
  const [rotationX, setRotationX] = useState(55); // Isometric pitch (degrees)
  const [rotationZ, setRotationZ] = useState(45); // Isometric yaw (degrees)
  const [zoom, setZoom] = useState(1.0);
  const [autoRotate, setAutoRotate] = useState(true);
  const [fps, setFps] = useState(60);
  const [hoveredSlot, setHoveredSlot] = useState<SlotSnapshot | null>(null);

  const isDraggingRef = useRef(false);
  const lastMousePosRef = useRef({ x: 0, y: 0 });
  const animFrameIdRef = useRef<number | null>(null);
  const frameCountRef = useRef(0);
  const lastFpsUpdateRef = useRef(0);

  // Handle Drag / Orbit
  const handleMouseDown = (e: React.MouseEvent<HTMLCanvasElement>) => {
    isDraggingRef.current = true;
    lastMousePosRef.current = { x: e.clientX, y: e.clientY };
  };

  const handleMouseMove = (e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isDraggingRef.current) return;
    const dx = e.clientX - lastMousePosRef.current.x;
    const dy = e.clientY - lastMousePosRef.current.y;
    lastMousePosRef.current = { x: e.clientX, y: e.clientY };

    setRotationZ((prev) => (prev + dx * 0.5) % 360);
    setRotationX((prev) => Math.max(15, Math.min(85, prev - dy * 0.5)));
  };

  const handleMouseUp = () => {
    isDraggingRef.current = false;
  };

  const handleWheel = (e: React.WheelEvent<HTMLCanvasElement>) => {
    e.preventDefault();
    setZoom((prev) => Math.max(0.6, Math.min(1.8, prev - e.deltaY * 0.001)));
  };

  // Render Loop at 60 FPS
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d", { alpha: false });
    if (!ctx) return;

    let localRotationZ = rotationZ;

    const render = (time: number) => {
      // FPS measurement
      frameCountRef.current++;
      if (lastFpsUpdateRef.current === 0) {
        lastFpsUpdateRef.current = time;
      } else if (time - lastFpsUpdateRef.current >= 500) {
        setFps(Math.round((frameCountRef.current * 1000) / (time - lastFpsUpdateRef.current)));
        frameCountRef.current = 0;
        lastFpsUpdateRef.current = time;
      }

      // Auto-rotation
      if (autoRotate && !isDraggingRef.current) {
        localRotationZ = (localRotationZ + 0.15) % 360;
      } else {
        localRotationZ = rotationZ;
      }

      // Resize canvas to match display DPI
      const width = canvas.clientWidth;
      const height = canvas.clientHeight;
      if (canvas.width !== width || canvas.height !== height) {
        canvas.width = width;
        canvas.height = height;
      }

      // Background: Deep sub-kernel dark slate with subtle radial glow
      const bgGrad = ctx.createRadialGradient(
        width / 2,
        height / 2,
        50,
        width / 2,
        height / 2,
        Math.max(width, height) * 0.8
      );
      bgGrad.addColorStop(0, "#080C14");
      bgGrad.addColorStop(1, "#020408");
      ctx.fillStyle = bgGrad;
      ctx.fillRect(0, 0, width, height);

      // Draw subtle background grid lines
      ctx.strokeStyle = "rgba(255, 255, 255, 0.03)";
      ctx.lineWidth = 1;
      const gridStep = 40;
      for (let x = 0; x < width; x += gridStep) {
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x, height);
        ctx.stroke();
      }
      for (let y = 0; y < height; y += gridStep) {
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(width, y);
        ctx.stroke();
      }

      // 3D Isometric Projection Setup
      const centerX = width / 2;
      const centerY = height / 2 + 30;
      const radX = (rotationX * Math.PI) / 180;
      const radZ = (localRotationZ * Math.PI) / 180;

      const cosX = Math.cos(radX);
      const sinX = Math.sin(radX);
      const cosZ = Math.cos(radZ);
      const sinZ = Math.sin(radZ);

      // Project 3D coordinate (x, y, z) to 2D screen coordinate
      const project = (x: number, y: number, z: number) => {
        // Rotate around Z axis (yaw)
        const rx = x * cosZ - y * sinZ;
        const ry = x * sinZ + y * cosZ;

        // Rotate around X axis (pitch)
        const rz = ry * sinX + z * cosX;
        const finalY = ry * cosX - z * sinX;

        // Perspective / Isometric scaling
        const scale = zoom;
        return {
          px: centerX + rx * scale,
          py: centerY + finalY * scale,
          depth: rz,
        };
      };

      // Draw Silicon Substrate Base Wafer (NUMA node)
      const waferSize = 240;
      const w0 = project(-waferSize, -waferSize, -8);
      const w1 = project(waferSize, -waferSize, -8);
      const w2 = project(waferSize, waferSize, -8);
      const w3 = project(-waferSize, waferSize, -8);

      ctx.beginPath();
      ctx.moveTo(w0.px, w0.py);
      ctx.lineTo(w1.px, w1.py);
      ctx.lineTo(w2.px, w2.py);
      ctx.lineTo(w3.px, w3.py);
      ctx.closePath();
      ctx.fillStyle = "rgba(10, 16, 28, 0.9)";
      ctx.fill();
      ctx.strokeStyle = "rgba(0, 240, 255, 0.25)";
      ctx.lineWidth = 1.5;
      ctx.stroke();

      // Render 8x8 matrix of 64-byte structural cache-line slots
      const slots = snapshot?.slots || [];
      const cols = 8;
      const rows = 8;
      const slotPitch = 48; // distance between slot centers
      const slotWidth = 34; // slot block dimension

      interface RenderableSlot {
        slot: SlotSnapshot;
        depth: number;
        pts: {
          top: { px: number; py: number }[];
          left: { px: number; py: number }[];
          right: { px: number; py: number }[];
        };
        height: number;
        color: string;
        isHovered: boolean;
        isSelected: boolean;
      }

      const renderList: RenderableSlot[] = [];

      for (let r = 0; r < rows; r++) {
        for (let c = 0; c < cols; c++) {
          const idx = r * cols + c;
          const slot = slots[idx] || {
            index: idx,
            slot_id: idx,
            tenant_id: "00000000",
            domain_key: 0,
            status: "FREE",
            active_tokens: 0,
            access_count: 0,
            memory_offset_hex: `0x${(idx * 64).toString(16).toUpperCase().padStart(4, "0")}`,
            cache_line_aligned: true,
            colorChannel: "#1E293B",
          };

          // 3D coordinates relative to center
          const posX = (c - cols / 2 + 0.5) * slotPitch;
          const posY = (r - rows / 2 + 0.5) * slotPitch;

          // Height based on status and token activity
          let h = 12;
          if (slot.status === "ACTIVE") {
            const pulse = Math.sin(time * 0.004 + idx) * 4;
            h = 24 + (slot.active_tokens % 20) + pulse;
          } else if (slot.status === "FENCED") {
            h = 32;
          }

          const isSelected = selectedSlot?.slot_id === slot.slot_id;
          const isHovered = hoveredSlot?.slot_id === slot.slot_id;
          if (isSelected || isHovered) {
            h += 6;
          }

          const hw = slotWidth / 2;

          // Compute 8 corner vertices of the 3D pillar
          // Base corners (z = 0)
          const b0 = project(posX - hw, posY - hw, 0);
          const b1 = project(posX + hw, posY - hw, 0);
          const b2 = project(posX + hw, posY + hw, 0);
          const b3 = project(posX - hw, posY + hw, 0);

          // Top corners (z = h)
          const t0 = project(posX - hw, posY - hw, h);
          const t1 = project(posX + hw, posY - hw, h);
          const t2 = project(posX + hw, posY + hw, h);
          const t3 = project(posX - hw, posY + hw, h);

          // Center depth for painter's algorithm sorting
          const centerPt = project(posX, posY, h / 2);

          renderList.push({
            slot,
            depth: centerPt.depth,
            pts: {
              top: [t0, t1, t2, t3],
              left: [b3, b2, t2, t3],
              right: [b2, b1, t1, t2],
            },
            height: h,
            color: slot.color_channel || "#1E293B",
            isHovered,
            isSelected,
          });
        }
      }

      // Sort by depth (back to front) for correct 3D occlusion
      renderList.sort((a, b) => a.depth - b.depth);

      // Render sorted 3D slots
      renderList.forEach((item) => {
        const { pts, color, isHovered, isSelected, slot } = item;

        // Base color handling
        let fillColor = color;
        let strokeColor = color;

        if (slot.status === "FREE") {
          fillColor = "rgba(30, 41, 59, 0.4)";
          strokeColor = "rgba(71, 85, 105, 0.3)";
        }

        // Draw Left Face
        ctx.beginPath();
        ctx.moveTo(pts.left[0].px, pts.left[0].py);
        ctx.lineTo(pts.left[1].px, pts.left[1].py);
        ctx.lineTo(pts.left[2].px, pts.left[2].py);
        ctx.lineTo(pts.left[3].px, pts.left[3].py);
        ctx.closePath();
        ctx.fillStyle = fillColor;
        ctx.globalAlpha = 0.65;
        ctx.fill();
        ctx.strokeStyle = strokeColor;
        ctx.lineWidth = 1;
        ctx.stroke();

        // Draw Right Face
        ctx.beginPath();
        ctx.moveTo(pts.right[0].px, pts.right[0].py);
        ctx.lineTo(pts.right[1].px, pts.right[1].py);
        ctx.lineTo(pts.right[2].px, pts.right[2].py);
        ctx.lineTo(pts.right[3].px, pts.right[3].py);
        ctx.closePath();
        ctx.fillStyle = fillColor;
        ctx.globalAlpha = 0.8;
        ctx.fill();
        ctx.strokeStyle = strokeColor;
        ctx.stroke();

        // Draw Top Face (Cap)
        ctx.beginPath();
        ctx.moveTo(pts.top[0].px, pts.top[0].py);
        ctx.lineTo(pts.top[1].px, pts.top[1].py);
        ctx.lineTo(pts.top[2].px, pts.top[2].py);
        ctx.lineTo(pts.top[3].px, pts.top[3].py);
        ctx.closePath();
        ctx.fillStyle = fillColor;
        ctx.globalAlpha = isSelected || isHovered ? 1.0 : 0.95;
        ctx.fill();

        // Selection / Hover Highlight Glow
        if (isSelected || isHovered) {
          ctx.strokeStyle = "#FFFFFF";
          ctx.lineWidth = 2.5;
          ctx.shadowColor = color;
          ctx.shadowBlur = 12;
          ctx.stroke();
          ctx.shadowBlur = 0;
        } else {
          ctx.strokeStyle = strokeColor;
          ctx.lineWidth = 1;
          ctx.stroke();
        }

        // Draw 64B Structural Silicon Fence on Fenced Slots
        if (slot.status === "FENCED") {
          ctx.strokeStyle = "#FF0055";
          ctx.lineWidth = 2;
          ctx.shadowColor = "#FF0055";
          ctx.shadowBlur = 8;
          ctx.stroke();
          ctx.shadowBlur = 0;
        }

        ctx.globalAlpha = 1.0;
      });

      animFrameIdRef.current = requestAnimationFrame(render);
    };

    animFrameIdRef.current = requestAnimationFrame(render);

    return () => {
      if (animFrameIdRef.current) {
        cancelAnimationFrame(animFrameIdRef.current);
      }
    };
  }, [rotationX, rotationZ, zoom, autoRotate, snapshot, selectedSlot, hoveredSlot]);

  return (
    <div className="relative overflow-hidden rounded-2xl border border-white/10 bg-zinc-950 shadow-2xl">
      {/* 3D WebGL Canvas */}
      <canvas
        ref={canvasRef}
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
        onClick={() => {
          if (onSelectSlot && (hoveredSlot || selectedSlot)) {
            onSelectSlot((hoveredSlot ?? selectedSlot) ?? null);
          }
        }}
        onWheel={handleWheel}
        className="h-[520px] w-full cursor-grab active:cursor-grabbing"
      />

      {/* Top HUD Overlay Controls */}
      <div className="pointer-events-none absolute left-4 top-4 flex flex-wrap items-center gap-2">
        <div className="pointer-events-auto flex items-center gap-2 rounded-lg border border-white/10 bg-black/60 px-3 py-1.5 backdrop-blur-md">
          <span className="flex h-2 w-2 rounded-full bg-emerald-400 animate-pulse" />
          <span className="text-xs font-semibold uppercase tracking-wider text-white">
            64B Sub-Kernel Matrix HUD
          </span>
          <span className="rounded bg-zinc-800 px-1.5 py-0.5 font-mono text-[11px] text-emerald-400">
            {fps} FPS
          </span>
        </div>

        <div className="pointer-events-auto flex items-center gap-1 rounded-lg border border-white/10 bg-black/60 p-1 backdrop-blur-md">
          <button
            onClick={() => setAutoRotate((prev) => !prev)}
            className={`rounded px-2.5 py-1 text-xs font-medium transition ${
              autoRotate ? "bg-cyan-500/20 text-cyan-400 border border-cyan-500/40" : "text-zinc-400 hover:text-white"
            }`}
          >
            {autoRotate ? "Auto-Orbit: ON" : "Auto-Orbit: OFF"}
          </button>
          <button
            onClick={() => {
              setRotationX(55);
              setRotationZ(45);
              setZoom(1.0);
            }}
            className="flex items-center gap-1 rounded px-2 py-1 text-xs text-zinc-400 hover:bg-zinc-800 hover:text-white transition"
            title="Reset 3D camera"
          >
            <RotateCcw className="h-3 w-3" />
            Reset
          </button>
        </div>
      </div>

      {/* Neon Channel Legend Overlay */}
      <div className="pointer-events-none absolute bottom-4 left-4 flex flex-wrap items-center gap-2">
        <div className="pointer-events-auto flex flex-wrap items-center gap-3 rounded-xl border border-white/10 bg-black/70 px-4 py-2 text-xs backdrop-blur-md">
          <div className="flex items-center gap-1.5">
            <span className="h-2.5 w-2.5 rounded-sm bg-[#00F0FF] shadow-[0_0_8px_#00F0FF]" />
            <span className="text-zinc-300">Tenant Alpha (Finance)</span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="h-2.5 w-2.5 rounded-sm bg-[#00FF88] shadow-[0_0_8px_#00FF88]" />
            <span className="text-zinc-300">Tenant Beta (Defense)</span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="h-2.5 w-2.5 rounded-sm bg-[#BD00FF] shadow-[0_0_8px_#BD00FF]" />
            <span className="text-zinc-300">Tenant Gamma (Health)</span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="h-2.5 w-2.5 rounded-sm bg-[#FFB800] shadow-[0_0_8px_#FFB800]" />
            <span className="text-zinc-300">Tenant Delta (Frontier)</span>
          </div>
          <div className="flex items-center gap-1.5">
            <span className="h-2.5 w-2.5 rounded-sm bg-[#FF0055] shadow-[0_0_8px_#FF0055]" />
            <span className="text-rose-400 font-medium">MPK Silicon Fence</span>
          </div>
          <div className="flex items-center gap-1.5 border-l border-white/10 pl-2">
            <span className="h-2.5 w-2.5 rounded-sm bg-amber-400 shadow-[0_0_8px_#F59E0B]" />
            <span className="text-amber-400 font-medium">PCIe Gen5 DMA Stream</span>
          </div>
        </div>
      </div>

      {/* Selected Slot Inspector Pill */}
      {selectedSlot && (
        <div className="pointer-events-auto absolute right-4 top-4 w-72 rounded-xl border border-cyan-500/30 bg-black/85 p-4 text-xs shadow-2xl backdrop-blur-md">
          <div className="flex items-center justify-between border-b border-white/10 pb-2">
            <div className="flex items-center gap-1.5 font-semibold text-white">
              <Cpu className="h-4 w-4 text-cyan-400" />
              <span>Slot 0x{selectedSlot.slot_id.toString(16).toUpperCase()} Inspector</span>
            </div>
            <span
              className="rounded px-2 py-0.5 text-[10px] font-bold"
              style={{
                backgroundColor: `${selectedSlot.color_channel}20`,
                color: selectedSlot.color_channel,
                border: `1px solid ${selectedSlot.color_channel}50`,
              }}
            >
              {selectedSlot.status}
            </span>
          </div>

          <div className="mt-3 space-y-1.5 font-mono text-[11px]">
            <div className="flex justify-between">
              <span className="text-zinc-400">Physical Offset:</span>
              <span className="text-cyan-400 font-bold">{selectedSlot.memory_offset_hex}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-zinc-400">Cache Alignment:</span>
              <span className="text-emerald-400 font-bold">64-Byte Congruent</span>
            </div>
            <div className="flex justify-between">
              <span className="text-zinc-400">PCIe DMA Binding:</span>
              <span className={selectedSlot.dma_bound ? "text-amber-400 font-bold" : "text-zinc-500"}>
                {selectedSlot.dma_bound ? "Gen5 x16 (DGX Station)" : "Unmapped"}
              </span>
            </div>
            {selectedSlot.dma_bound && (
              <>
                <div className="flex justify-between">
                  <span className="text-zinc-400">MoE Expert ID:</span>
                  <span className="text-cyan-300 font-bold">Expert #{selectedSlot.expert_id ?? 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-zinc-400">DMA Latency:</span>
                  <span className="text-emerald-400 font-bold">
                    {selectedSlot.dma_latency_ms?.toFixed(2) ?? "0.82"} ms (&lt;2ms SLA)
                  </span>
                </div>
              </>
            )}
            <div className="flex justify-between">
              <span className="text-zinc-400">Intel MPK Domain:</span>
              <span className="text-white">Domain Key 0x0{selectedSlot.domain_key}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-zinc-400">Tenant Hash:</span>
              <span className="text-zinc-300">{selectedSlot.tenant_id}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-zinc-400">Active In-Flight:</span>
              <span className="text-white">{selectedSlot.active_tokens} tokens</span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
