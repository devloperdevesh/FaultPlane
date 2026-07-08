# Runtime

This document describes how AgentMesh behaves during execution.

It focuses on runtime behavior rather than implementation details. The objective is to explain how requests move through the system, how execution state evolves, and how recovery decisions are made.

---

# Runtime Overview

AgentMesh consists of two logical layers.

```text
               Control Plane
        ┌────────────────────────┐
        │ Gateway                │
        │ Routing                │
        │ Recovery               │
        │ Telemetry              │
        └──────────┬─────────────┘
                   │
                   ▼
               Data Plane
        ┌────────────────────────┐
        │ Primary Worker         │
        │ Fallback Worker        │
        │ Application Runtime    │
        └────────────────────────┘
```

The control plane coordinates execution.

The data plane performs application work.

---

# Runtime Lifecycle

A request passes through a fixed sequence of stages.

```text
Receive Request
        │
        ▼
Validate
        │
        ▼
Select Worker
        │
        ▼
Forward Request
        │
        ▼
Receive Response
        │
        ▼
Update Checkpoint
        │
        ▼
Return Response
```

Recovery is introduced only if execution cannot continue normally.

---

# Gateway Lifecycle

The gateway owns request coordination.

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
Accept Requests
   │
   ▼
Graceful Shutdown
```

The gateway never executes workload logic.

---

# Worker Lifecycle

Workers execute user workloads independently from the gateway.

```text
Start Worker
      │
      ▼
Receive Request
      │
      ▼
Execute Step
      │
      ▼
Checkpoint Update
      │
      ▼
Return Response
```

Workers are expected to be replaceable.

The control plane should not depend on any individual worker instance.

---

# Request Processing

Incoming requests follow the same processing pipeline.

```text
Client
   │
   ▼
Gateway
   │
   ▼
Routing
   │
   ▼
Worker
   │
   ▼
Response
```

The gateway maintains request metadata required for recovery.

---

# Runtime State

AgentMesh distinguishes between application state and infrastructure state.

| State | Owner |
|--------|-------|
| Request metadata | Gateway |
| Routing information | Control Plane |
| Checkpoint | Storage |
| Business data | Application |
| Execution context | Worker |

Only checkpoint information is shared across recovery events.

---

# Context Flow

Execution context evolves as work progresses.

```text
Request
    │
    ▼
Worker
    │
    ▼
Checkpoint
    │
    ▼
Storage
    │
    ▼
Recovery
```

Context should be captured after successful execution steps.

Failed operations must not overwrite the last valid checkpoint.

---

# Recovery Trigger

Recovery is initiated only after infrastructure failures.

Typical triggers include:

- worker unavailable
- connection timeout
- HTTP 5xx
- network interruption

Application-specific failures remain outside the runtime recovery model.

---

# Recovery Pipeline

```text
Worker Failure
      │
      ▼
Detect Failure
      │
      ▼
Lookup Checkpoint
      │
      ▼
Select Fallback
      │
      ▼
Restore Context
      │
      ▼
Resume Execution
```

Recovery is deterministic.

The same checkpoint should always produce the same recovery decision.

---

# Runtime Guarantees

The runtime is designed around a small number of guarantees.

| Guarantee | Description |
|-----------|-------------|
| Stateless Gateway | Execution progress is not stored in gateway memory. |
| Explicit Recovery | Recovery occurs only after validated failures. |
| Replaceable Workers | Workers may be restarted without changing gateway behavior. |
| Observable Execution | Recovery decisions are visible through telemetry. |

These guarantees define expected runtime behavior across all deployments.
---

# Failure Handling

Failures are classified before recovery decisions are made.

The runtime distinguishes between infrastructure failures and application failures.

| Failure Type | Recoverable | Example |
|--------------|-------------|---------|
| Worker crash | Yes | Process termination |
| HTTP 5xx | Yes | Upstream unavailable |
| Network timeout | Yes | Connection timeout |
| Connection refused | Yes | Worker offline |
| Invalid request | No | Malformed payload |
| Business logic error | No | Application exception |

Only infrastructure failures trigger the recovery pipeline.

---

# Health Checks

Workers expose lightweight health endpoints.

Typical health states include:

| State | Description |
|--------|-------------|
| Healthy | Worker accepts requests |
| Degraded | Worker responds with increased latency |
| Unavailable | Worker cannot process requests |

Health checks influence routing decisions but never modify application state.

---

# Concurrency Model

AgentMesh does not execute workloads directly.

Its concurrency model focuses on coordinating request routing and checkpoint updates.

```text
Incoming Requests
        │
        ▼
Gateway
        │
        ├─────────────┐
        ▼             ▼
 Worker A       Worker B
        │             │
        └──────┬──────┘
               ▼
      Checkpoint Store
```

Routing decisions should remain independent across concurrent requests.

---

# Scheduling

AgentMesh intentionally avoids implementing a workload scheduler.

Scheduling remains the responsibility of external systems such as:

- Kubernetes
- Container runtimes
- Orchestration frameworks

The gateway selects an execution target but does not decide when workloads should run.

---

# Resource Management

The control plane should maintain a small resource footprint.

Operational priorities include:

- low memory usage
- predictable CPU utilization
- minimal allocations
- fast startup
- graceful shutdown

The gateway should remain lightweight regardless of workload complexity.

---

# Graceful Shutdown

Shutdown follows a controlled sequence.

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
Terminate Process
```

New requests should not be accepted after shutdown begins.

---

# Timeout Strategy

Timeout values should be configured explicitly.

Different stages of the runtime may require different timeout thresholds.

Examples include:

- request timeout
- worker timeout
- checkpoint timeout
- telemetry export timeout

Using explicit timeouts improves operational predictability.

---

# Retry Strategy

AgentMesh does not perform unlimited retries.

Recovery attempts should remain bounded.

Typical strategy:

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
Done  Return Error
```

Repeated failures should be surfaced to operators rather than retried indefinitely.

---

# Runtime Constraints

Current implementation intentionally limits scope.

Current constraints include:

- single gateway instance
- in-memory checkpoints
- HTTP transport
- local development environment

These constraints simplify implementation while the recovery model matures.

---

# Performance Considerations

Runtime optimizations focus on predictable behavior.

Key areas include:

- routing latency
- checkpoint overhead
- allocation reduction
- telemetry efficiency
- startup time

Performance work should be driven by benchmark data rather than assumptions.

---

# Observability

Runtime events should be observable through telemetry.

Examples include:

| Event | Description |
|--------|-------------|
| Gateway startup | Control plane initialization |
| Request received | Incoming workload |
| Worker selected | Routing decision |
| Checkpoint updated | Successful execution step |
| Recovery initiated | Infrastructure failure detected |
| Recovery completed | Execution resumed |
| Shutdown complete | Runtime termination |

Operational visibility is considered a core runtime capability.

---

# Operational Recommendations

Recommended deployment practices include:

- monitor worker health
- validate recovery regularly
- export telemetry centrally
- minimize gateway dependencies
- keep checkpoint payloads small
- test failure scenarios continuously

Reliable recovery depends on operational discipline as much as implementation quality.

---

# Runtime Summary

The runtime is designed around a simple execution model.

```text
Request
   │
   ▼
Gateway
   │
   ▼
Worker
   │
Checkpoint
   │
   ▼
Storage
   │
   ▼
Recovery
   │
   ▼
Continue
```

The gateway coordinates execution.

Workers perform application logic.

Storage preserves execution state.

Recovery restores progress when infrastructure failures occur.

Each component remains focused on a single responsibility.

---