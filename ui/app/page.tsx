"use client";

import React from "react";
import {
  Users,
  Network,
  Activity,
  Database,
  DollarSign,
  CheckCircle2,
  AlertTriangle,
  AlertCircle,
  RefreshCw,
  Radio,
} from "lucide-react";

export default function OperationsDashboard() {
  return (
    <div className="p-6 max-w-[1600px] w-full mx-auto space-y-6">
      {/* 📡 THE INTEGRATED TOP DATA SUB-BAR */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center bg-[#0a0d1a] border border-slate-900 rounded-xl p-4 gap-4 shadow-lg">
        <div>
          <h2 className="text-sm font-bold text-white tracking-tight flex items-center gap-2">
            Overview{" "}
            <span className="text-xs font-mono font-normal text-slate-500">
              Real-time control plane diagnostics
            </span>
          </h2>
        </div>
        <div className="flex items-center gap-3 font-mono text-xs">
          <div className="bg-slate-950 border border-slate-900 px-3 py-1.5 rounded-lg text-slate-400">
            <span className="text-slate-600 font-semibold mr-1">
              Timeframe:
            </span>{" "}
            Last 30 Minutes
          </div>
          <div className="flex items-center gap-2 bg-slate-950 border border-slate-900 px-3 py-1.5 rounded-lg text-emerald-400">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
            <span className="uppercase text-[10px] font-bold tracking-wider">
              Auto Refresh
            </span>
          </div>
        </div>
      </div>

      {/* 📊 CORE TELEMETRY METRICS LAYER COUNTERS */}
      <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
        <StatCard
          title="Active Agents"
          value="248"
          change="▲ 12% vs last 30m"
          changeColor="text-emerald-500"
          icon={<Users className="text-purple-400/80" size={16} />}
        />
        <StatCard
          title="Active Workflows"
          value="1240"
          change="▲ 8% vs last 30m"
          changeColor="text-emerald-500"
          icon={<Network className="text-blue-400/80" size={16} />}
        />
        <StatCard
          title="Recovery Events"
          value="97"
          change="▼ -3% vs last 30m"
          changeColor="text-rose-500"
          icon={<Activity className="text-orange-400/80" size={16} />}
        />
        <StatCard
          title="Checkpoints Saved"
          value="3.42K"
          change="▲ 18% vs last 30m"
          changeColor="text-emerald-500"
          icon={<Database className="text-amber-400/80" size={16} />}
        />
        <StatCard
          title="Compute Saved"
          value="$12,430"
          change="▲ 24% vs last 30m"
          changeColor="text-emerald-500"
          icon={<DollarSign className="text-emerald-400/80" size={16} />}
        />
      </section>

      {/* ⚡ ACTIVE STEPS TIMELINE & LATENCY GAUGES */}
      <section className="grid gap-6 lg:grid-cols-3">
        {/* Failover Process Stepper */}
        <div className="lg:col-span-2 bg-[#0c101f] border border-slate-900 rounded-xl p-5 shadow-xl">
          <div className="flex justify-between items-center mb-5">
            <h3 className="text-xs font-bold font-mono tracking-widest text-slate-400 uppercase flex items-center gap-2">
              <Radio size={14} className="text-blue-500 animate-pulse" />{" "}
              Recovery Timeline
            </h3>
            <span className="text-[10px] font-mono text-emerald-400 bg-emerald-950/20 px-2 py-0.5 rounded border border-emerald-900/30">
              ● Live Stream
            </span>
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
            <TimelineStep
              time="10:24:31"
              title="Checkpoint"
              detail="Step 12"
              status="success"
              icon={<CheckCircle2 size={12} />}
            />
            <TimelineStep
              time="10:24:35"
              title="Worker Failure"
              detail="worker-02"
              status="warn"
              icon={<AlertTriangle size={12} />}
            />
            <TimelineStep
              time="10:24:36"
              title="Failure Detected"
              detail="TimeoutError"
              status="error"
              icon={<AlertCircle size={12} />}
            />
            <TimelineStep
              time="10:24:37"
              title="Recovery Init"
              detail="Step 12"
              status="process"
              icon={<RefreshCw className="animate-spin" size={12} />}
            />
            <TimelineStep
              time="10:24:38"
              title="Resumed"
              detail="Step 13"
              status="success"
              icon={<CheckCircle2 size={12} />}
            />
            <TimelineStep
              time="10:24:42"
              title="Completed"
              detail="Success"
              status="final"
              icon={<CheckCircle2 size={12} />}
            />
          </div>
        </div>

        {/* Latency Quantile Card */}
        <div className="bg-[#0c101f] border border-slate-900 rounded-xl p-5 shadow-xl flex flex-col justify-between">
          <div className="flex justify-between items-center font-mono text-xs text-slate-500">
            <span>Latency (P95)</span>
            <span>Live</span>
          </div>
          <div className="flex items-baseline gap-2 my-2">
            <span className="text-3xl font-extrabold text-slate-100 font-mono">
              2.3ms
            </span>
            <span className="text-[10px] font-mono font-semibold text-emerald-400 bg-emerald-950/30 px-1.5 py-0.5 rounded border border-emerald-900/30">
              ▲ 18% Fast
            </span>
          </div>
          <p className="text-[11px] text-slate-500 leading-normal">
            Low-overhead network failover routing operating seamlessly
            underneath the active request pipelines.
          </p>
        </div>
      </section>

      {/* 🏛️ BARE-METAL CHECKLISTS & HISTORY ROWS */}
      <section className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {/* Pipelines Log History Table */}
        <div className="lg:col-span-2 bg-[#0c101f] border border-slate-900 rounded-xl p-5 shadow-xl overflow-hidden">
          <h3 className="text-xs font-bold font-mono tracking-widest text-slate-400 uppercase mb-4">
            Recent Workflows
          </h3>
          <div className="overflow-x-auto">
            <table className="w-full font-mono text-xs text-left border-collapse">
              <thead>
                <tr className="border-b border-slate-900 text-slate-500 font-bold">
                  <th className="pb-2.5 uppercase text-[10px]">Workflow</th>
                  <th className="pb-2.5 uppercase text-[10px]">Agent</th>
                  <th className="pb-2.5 uppercase text-[10px]">Status</th>
                  <th className="pb-2.5 uppercase text-[10px] text-right">
                    Recoveries
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-900/40">
                <TableRow
                  name="refund-processing"
                  agent="customer-support-agent-42"
                  status="Completed"
                  color="text-emerald-400"
                  count="1"
                />
                <TableRow
                  name="data-analysis"
                  agent="research-agent-07"
                  status="Running"
                  color="text-blue-400 animate-pulse"
                  count="0"
                />
                <TableRow
                  name="report-generation"
                  agent="report-agent-03"
                  status="Completed"
                  color="text-emerald-400"
                  count="0"
                />
                <TableRow
                  name="email-summarization"
                  agent="email-agent-11"
                  status="Failed"
                  color="text-rose-400"
                  count="2"
                />
              </tbody>
            </table>
          </div>
        </div>

        {/* Infrastructure Checklist Panel */}
        <div className="bg-[#0c101f] border border-slate-900 rounded-xl p-5 shadow-xl">
          <h3 className="text-xs font-bold font-mono tracking-widest text-slate-400 uppercase mb-4">
            System Health Status
          </h3>
          <div className="space-y-3 font-mono text-xs">
            <HealthItem label="Control Plane Daemon" value="100%" />
            <HealthItem label="Data Plane Core" value="100%" />
            <HealthItem label="Lock-Free Checkpoint Store" value="100%" />
            <HealthItem label="eBPF Event Stream" value="99.9%" />
            <HealthItem label="Telemetry Pipeline" value="99.8%" />
          </div>
        </div>
      </section>
    </div>
  );
}

// Inline Sub-Components
function StatCard({
  title,
  value,
  change,
  changeColor,
  icon,
}: {
  title: string;
  value: string;
  change: string;
  changeColor: string;
  icon: React.ReactNode;
}) {
  return (
    <div className="bg-[#0c101f] border border-slate-900 rounded-xl p-4 shadow-md">
      <div className="flex justify-between items-center mb-2">
        <span className="text-[10px] font-bold font-mono uppercase tracking-widest text-slate-500">
          {title}
        </span>
        <div className="p-1.5 bg-slate-950 border border-slate-900 rounded-lg">
          {icon}
        </div>
      </div>
      <h4 className="text-xl font-extrabold font-mono text-slate-100 tracking-tight">
        {value}
      </h4>
      <p className={`text-[10px] font-mono font-medium mt-1 ${changeColor}`}>
        {change}
      </p>
    </div>
  );
}

function TimelineStep({
  time,
  title,
  detail,
  status,
  icon,
}: {
  time: string;
  title: string;
  detail: string;
  status: string;
  icon: React.ReactNode;
}) {
  const statusColors =
    status === "success"
      ? "bg-emerald-950/30 border-emerald-500/50 text-emerald-400"
      : status === "warn"
        ? "bg-amber-950/30 border-amber-500/50 text-amber-400"
        : status === "error"
          ? "bg-rose-950/30 border-rose-500/50 text-rose-400"
          : status === "process"
            ? "bg-blue-950/30 border-blue-500/50 text-blue-400"
            : "bg-slate-950 border-slate-800 text-slate-500";

  return (
    <div className="flex flex-col items-center text-center">
      <span className="text-[9px] font-mono text-slate-600 mb-1">{time}</span>
      <div
        className={`w-6 h-6 rounded-full flex items-center justify-center border font-mono ${statusColors}`}
      >
        {icon}
      </div>
      <p className="text-[10px] font-bold text-slate-200 mt-2 tracking-tight leading-tight">
        {title}
      </p>
      <span className="text-[9px] font-mono text-slate-500 mt-0.5">
        {detail}
      </span>
    </div>
  );
}

function TableRow({
  name,
  agent,
  status,
  color,
  count,
}: {
  name: string;
  agent: string;
  status: string;
  color: string;
  count: string;
}) {
  return (
    <tr className="border-b border-slate-900/20 text-slate-300 hover:bg-slate-950/20 transition-colors">
      <td className="py-2.5 font-semibold text-slate-200">{name}</td>
      <td className="py-2.5 text-slate-400">{agent}</td>
      <td className={`py-2.5 font-bold ${color}`}>{status}</td>
      <td className="py-2.5 font-bold text-right text-slate-500">{count}</td>
    </tr>
  );
}

function HealthItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="space-y-2">
      <div className="flex justify-between items-center">
        <span className="text-slate-400">{label}</span>

        <span className="text-emerald-400 font-bold">{value}</span>
      </div>

      <div className="w-full h-1.5 bg-slate-950 rounded-full overflow-hidden">
        <div
          className="bg-emerald-500 h-full rounded-full"
          style={{ width: value }}
        />
      </div>
    </div>
  );
}