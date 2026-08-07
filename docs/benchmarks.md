# Benchmarks

This document describes the benchmarking methodology used by FaultPlane.

The purpose of benchmarking is to measure infrastructure behavior, runtime overhead, and system characteristics under controlled conditions.

FaultPlane benchmarks focus on understanding:

- routing overhead
- recovery performance
- state management cost
- runtime resource usage
- scalability characteristics

Benchmark results should always include environment details and reproducible methodology.

---

# Benchmark Philosophy

Benchmarks exist to support engineering decisions, not marketing claims.

Performance improvements should be validated through:

- measurable experiments
- reproducible workloads
- consistent environments
- transparent reporting

A single benchmark number does not represent complete system performance.

---

# Objectives

The benchmark suite aims to answer:

- What latency overhead does FaultPlane introduce?
- How quickly can failures be detected and recovered?
- How much memory does the runtime consume?
- How does throughput behave under concurrent workloads?
- What is the cost of maintaining runtime state?
- How does observability affect runtime performance?

---

# Scope

FaultPlane benchmarks focus on infrastructure components.

## Included

- gateway routing latency
- recovery latency
- runtime state operations
- storage performance
- request throughput
- memory allocation
- CPU utilization
- telemetry overhead

---

## Excluded

The benchmark suite does not measure:

- LLM inference quality
- model accuracy
- prompt execution performance
- application business logic
- vector database performance
- external API latency

These depend on external systems and workloads.

---

# Benchmark Environment

Every published benchmark should document the execution environment.

Example:

| Component | Configuration |
|---|---|
| CPU | AMD Ryzen 9 7900X |
| Memory | 32GB DDR5 |
| Operating System | Ubuntu 24.04 LTS |
| Go Version | 1.23+ |
| Docker Version | Latest Stable |
| FaultPlane Version | Release / Git SHA |
| Runtime Configuration | Documented |

Hardware and software details are required for reproducibility.

---

# Benchmark Workflow

Benchmarks should execute using a controlled process.

```
Environment Setup

        ↓

Configuration Loading

        ↓

Warm-up Phase

        ↓

Benchmark Execution

        ↓

Metric Collection

        ↓

Result Analysis
```

Warm-up iterations should not be included in final measurements.

---

# Workload Categories

| Workload | Purpose |
|---|---|
| Single Request | Measure baseline latency |
| Concurrent Requests | Measure throughput behavior |
| Failure Recovery | Measure recovery overhead |
| State Writes | Measure checkpoint performance |
| State Reads | Measure lookup performance |
| Long Running Sessions | Evaluate runtime stability |

Additional workloads may be introduced as the runtime evolves.

---

# Metrics

FaultPlane tracks multiple performance indicators.

| Metric | Description |
|---|---|
| Request Latency | Time required to process a request |
| Recovery Latency | Time required to restore execution |
| Throughput | Requests handled per second |
| Memory Allocation | Heap allocation cost |
| Memory Usage | Runtime memory footprint |
| CPU Usage | Processing overhead |
| Error Rate | Reliability under workload |

---

# Benchmark Categories

## Gateway Routing

Measures the overhead introduced by the routing layer.

Flow:

```
Client

  ↓

FaultPlane Gateway

  ↓

Worker Runtime
```

Metrics:

- routing latency
- request throughput
- CPU overhead

---

# Recovery Performance

Measures failure handling behavior.

Flow:

```
Failure Event

      ↓

Detection

      ↓

State Lookup

      ↓

Recovery Decision

      ↓

Execution Resume
```

Metrics:

- detection time
- recovery duration
- successful recovery rate

---

# Runtime State

Measures state management performance.

Flow:

```
Execution State

        ↓

Checkpoint Operation

        ↓

Storage Layer
```

Metrics:

- write latency
- read latency
- storage overhead
- serialization cost

---

# Telemetry Overhead

Measures the performance impact of observability systems.

Evaluates:

- logging overhead
- metric collection cost
- trace generation impact
- export latency

Observability should provide visibility without significantly affecting runtime behavior.

---

# Benchmark Reporting

Every benchmark report should include:

## Environment

- hardware information
- software versions
- runtime configuration

## Workload

- request pattern
- concurrency level
- dataset size
- execution duration

## Results

- raw measurements
- summarized metrics
- comparison results

## Interpretation

Explain:

- what changed
- why it changed
- possible limitations

---

# Reproducibility Guidelines

Benchmarks should be repeatable.

Recommended practices:

- use fixed configurations
- document commands
- isolate test environments
- version benchmark code
- avoid hidden variables

Example:

```bash
go test -bench=. ./...
```

Benchmark methodology changes should be documented.

---

# Performance Priorities

FaultPlane optimization focuses on predictable production behavior.

Priority areas:

1. Recovery latency
2. Tail latency consistency
3. Memory efficiency
4. Runtime stability
5. Throughput scalability
6. Operational simplicity

Average latency alone is not sufficient for evaluating distributed infrastructure.

---

# Future Benchmark Work

Planned improvements:

- distributed recovery benchmarks
- Kubernetes deployment testing
- storage backend comparisons
- multi-node evaluation
- telemetry scaling analysis
- automated benchmark pipelines
- regression tracking

---

# Benchmark Lifecycle

```
Workload Definition

        ↓

Benchmark Execution

        ↓

Metric Collection

        ↓

Performance Analysis

        ↓

Engineering Decision
```

---

# Conclusion

Benchmarking helps FaultPlane evolve through measurable engineering decisions.

The goal is not only faster systems, but:

- predictable behavior
- reliable recovery
- efficient resource usage
- operational confidence

Performance improvements should always be supported by transparent measurements.
