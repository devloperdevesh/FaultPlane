# Runtime

This document describes how FaultPlane behaves during execution.

It explains runtime lifecycle, request flow, state ownership, failure handling, and recovery behavior.

Unlike implementation documentation, this document focuses on **runtime semantics** — how components interact during normal execution and infrastructure failures.

---

# Runtime Overview

FaultPlane consists of two logical layers:

```text
                 Control Plane

        ┌────────────────────────┐
        │ Gateway                │
        │ Routing Engine         │
        │ Recovery Manager       │
        │ Telemetry              │
        └──────────┬─────────────┘
                   │
                   ▼

                 Data Plane

        ┌────────────────────────┐
        │ Worker Runtime         │
        │ Application Execution  │
        │ External Services      │
        └────────────────────────┘
```

The control plane coordinates traffic and recovery.

The data plane performs workload execution.

---

# Runtime Lifecycle

Every request follows a predictable execution lifecycle.

```text
Receive Request

        │

        ▼

Validate Metadata

        │

        ▼

Select Runtime

        │

        ▼

Forward Request

        │

        ▼

Execute Work

        │

        ▼

Update Checkpoint

        │

        ▼

Return Response
```

Recovery is triggered only when execution cannot continue normally.

---

# Gateway Lifecycle

The gateway manages request coordination.

Lifecycle:

```text
Start

 │

 ▼

Load Configuration

 │

 ▼

Initialize Storage

 │

 ▼

Initialize Telemetry

 │

 ▼

Accept Traffic

 │

 ▼

Graceful Shutdown
```

The gateway does not execute application workloads.

Its responsibility is limited to:

- routing
- failure detection
- recovery coordination
- telemetry generation

---

# Worker Lifecycle

Workers execute workload operations independently.

```text
Start Worker

      │

      ▼

Receive Request

      │

      ▼

Execute Operation

      │

      ▼

Create Checkpoint

      │

      ▼

Return Result
```

Workers are replaceable execution targets.

The control plane should not depend on any specific worker instance.

---

# Request Processing Flow

Normal request execution:

```text
Client

   │

   ▼

FaultPlane Gateway

   │

   ▼

Routing Engine

   │

   ▼

Worker Runtime

   │

   ▼

Response
```

The gateway maintains only the metadata required for routing and recovery.

---

# Runtime State Model

FaultPlane separates infrastructure state from application state.

| State | Owner |
|---|---|
| Request Metadata | Gateway |
| Routing Decision | Control Plane |
| Execution State | Checkpoint System |
| Application Data | Workload |
| Runtime Health | Telemetry Layer |

Only recovery-related state is shared across runtime failures.

---

# Execution Context Flow

Execution state evolves through checkpoints.

```text
Request

   │

   ▼

Worker Runtime

   │

   ▼

Checkpoint Manager

   │

   ▼

Storage Backend

   │

   ▼

Recovery Runtime
```

Checkpoints should only represent successful execution progress.

Failed operations must not overwrite valid recovery state.

---

# Recovery Trigger

Recovery is initiated only for infrastructure failures.

Examples:

- worker crash
- connection timeout
- unavailable runtime
- network interruption
- upstream failure

Application-level failures remain outside the recovery model.

---

# Recovery Pipeline

Recovery follows a deterministic sequence.

```text
Runtime Failure

        │

        ▼

Failure Detection

        │

        ▼

Locate Checkpoint

        │

        ▼

Select Runtime

        │

        ▼

Restore State

        │

        ▼

Resume Execution
```

The same checkpoint should produce consistent recovery behavior.

---

# Runtime Guarantees

FaultPlane provides several runtime guarantees.

| Guarantee | Description |
|---|---|
| Stateless Gateway | Execution state does not depend on gateway memory. |
| Replaceable Workers | Workers can restart without losing progress. |
| Explicit Recovery | Recovery happens through defined workflows. |
| Observable Execution | Runtime events are measurable. |

These guarantees define expected behavior across deployments.

---

# Failure Classification

FaultPlane separates infrastructure failures from application failures.

| Failure Type | Recoverable | Example |
|---|---|---|
| Worker crash | Yes | Process termination |
| Container failure | Yes | Runtime restart |
| Timeout | Yes | Network delay |
| Connection failure | Yes | Worker unavailable |
| Invalid request | No | Bad input |
| Business error | No | Application exception |

Only infrastructure failures enter recovery workflows.

---

# Health Model

Workers expose lightweight health information.

Example states:

| State | Description |
|---|---|
| Healthy | Accepting workload |
| Degraded | Reduced reliability |
| Unavailable | Cannot process requests |

Health information influences routing decisions.

It does not modify application state.

---

# Concurrency Model

FaultPlane coordinates concurrent execution without directly managing application threads.

Example:

```text
Incoming Requests

        │

        ▼

Gateway

   ┌────┴────┐

   ▼         ▼

Worker A   Worker B

   │         │

   └────┬────┘

        ▼

Checkpoint Storage
```

Routing decisions remain isolated between concurrent requests.

---

# Scheduling Responsibility

FaultPlane does not implement workload scheduling.

Scheduling remains the responsibility of external systems:

- Kubernetes
- container platforms
- orchestration systems

FaultPlane selects execution targets but does not decide when workloads should execute.

---

# Resource Management

The runtime is designed to maintain a lightweight operational footprint.

Important properties:

- low memory usage
- predictable CPU usage
- minimal allocations
- fast startup
- graceful shutdown

The gateway should remain efficient regardless of workload complexity.

---

# Graceful Shutdown

Shutdown follows a controlled lifecycle.

```text
Stop Accepting Requests

        │

        ▼

Complete Active Requests

        │

        ▼

Flush Telemetry

        │

        ▼

Release Resources

        │

        ▼

Terminate Runtime
```

New traffic should not enter the system after shutdown begins.

---

# Timeout Strategy

Timeouts should be explicit and configurable.

Important timeout categories:

- request timeout
- worker timeout
- checkpoint timeout
- telemetry timeout

Explicit timeout management improves operational predictability.

---

# Retry Strategy

FaultPlane avoids unlimited retries.

Recovery attempts should remain controlled.

Flow:

```text
Request

   │

   ▼

Worker

   │

Failure

   │

   ▼

Recovery

   │

Success?

 ┌──┴──┐

 │     │

Yes    No

 │     │

 ▼     ▼

Done  Error
```

Repeated failures should become visible through telemetry.

---

# Current Runtime Constraints

Current implementation focuses on establishing the recovery foundation.

Current limitations:

- single gateway instance
- in-memory checkpoint storage
- HTTP transport
- local development focus

These constraints allow the recovery model to mature before introducing distributed complexity.

---

# Performance Considerations

Runtime optimization focuses on predictable behavior.

Primary areas:

- routing latency
- checkpoint overhead
- memory allocation
- telemetry cost
- startup performance

Optimization decisions should be validated through benchmarks.

---

# Observability

Runtime events should be visible through telemetry.

Examples:

| Event | Description |
|---|---|
| Gateway Started | Control plane initialization |
| Request Received | Incoming traffic |
| Runtime Selected | Routing decision |
| Checkpoint Created | Execution progress saved |
| Recovery Started | Failure handling initiated |
| Recovery Completed | Execution restored |
| Shutdown Complete | Runtime terminated |

Observability is a core runtime capability.

---

# Operational Recommendations

Recommended practices:

- monitor worker health
- validate recovery paths
- collect runtime metrics
- test failure scenarios
- keep checkpoints efficient
- review telemetry regularly

Reliable recovery requires both correct architecture and disciplined operations.

---

# Runtime Summary

FaultPlane follows a simple execution model:

```text
Request

   │

   ▼

Gateway

   │

   ▼

Worker Runtime

   │

   ▼

Checkpoint

   │

   ▼

Storage

   │

   ▼

Recovery

   │

   ▼

Continue Execution
```

Responsibilities remain separated:

- Gateway coordinates traffic.
- Workers execute workloads.
- Storage preserves execution state.
- Recovery restores progress.

This separation allows FaultPlane to provide resilience for long-running AI infrastructure while keeping runtime behavior predictable and observable.
