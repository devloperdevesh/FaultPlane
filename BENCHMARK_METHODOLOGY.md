# FaultPlane Infrastructure Performance and Benchmark Methodology

This document defines the benchmarking methodology used to evaluate FaultPlane transport efficiency under controlled enterprise workloads.

## 1. Simulation Testing Matrix

### Kernel Host Harness

Buildroot-generated Linux Kernel v6.6+ environments are used with multi-node QEMU and `virtme-ng` hardware-accelerated sandboxes where applicable.

### Telemetry Instrumentation

Kernel-level instrumentation uses `kprobe`, `bpftrace`, and eBPF mechanisms to collect packet-processing and execution telemetry.

### Volumetric Load Generation

Distributed high-throughput workload replay is used to evaluate FaultPlane under large numbers of concurrent autonomous-agent network routes.

Workloads should document:

- Concurrent connection count
- Packet rate
- Payload characteristics
- Traffic duration
- Failover frequency
- CPU allocation
- Memory allocation
- Kernel configuration

## 2. Performance Metrics

The following values represent target performance thresholds to be validated through reproducible benchmark runs:

| Metric | Target |
| --- | ---: |
| P95 Latency | 1.42 ms |
| P99 Tail Latency | 1.84 ms |
| CPU Processing Overhead | < 0.04% |

These targets must not be interpreted as measured production results unless accompanied by reproducible benchmark evidence.

## 3. Benchmark Reproducibility

Every reported benchmark result should include:

1. Hardware configuration
2. Kernel version
3. FaultPlane commit or release version
4. Workload configuration
5. Number of samples
6. Measurement tooling
7. Benchmark command
8. Raw or summarized results
9. Relevant system configuration

Benchmark results should be reproducible by an independent reviewer using the documented environment and methodology.

## 4. Production Measurements

Production telemetry must be kept separate from benchmark targets.

Production measurements should be recorded in the appropriate deployment telemetry and audit-log infrastructure and should include sufficient context to distinguish observed results from predefined performance targets.