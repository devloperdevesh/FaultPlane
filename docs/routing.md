# Routing

This document describes how AgentMesh routes requests between execution workers.

Routing is responsible for selecting an appropriate execution target while remaining independent from application logic.

The routing layer does not execute workloads, interpret prompts, or modify execution state. Its responsibility is limited to request forwarding and recovery coordination.

---

# Overview

Every incoming request passes through the routing subsystem before reaching a worker.

```text
Client
   │
   ▼
Gateway
   │
   ▼
Routing Engine
   │
   ▼
Worker
```

The routing engine determines where a request should execute based on infrastructure state rather than application behavior.

---

# Objectives

The routing layer is designed around the following goals.

| Goal | Description |
|------|-------------|
| Availability | Route requests to healthy workers. |
| Recovery | Redirect execution after infrastructure failures. |
| Simplicity | Keep routing deterministic and predictable. |
| Low Latency | Introduce minimal routing overhead. |
| Observability | Record routing decisions through telemetry. |

---

# Responsibilities

The routing subsystem is responsible for:

- selecting execution targets
- evaluating worker availability
- coordinating recovery
- forwarding requests
- exposing routing metrics

The routing subsystem is not responsible for:

- scheduling workloads
- load balancing clusters
- executing business logic
- model selection
- workflow orchestration

---

# Routing Pipeline

Every request follows the same routing sequence.

```text
Receive Request
        │
        ▼
Validate Metadata
        │
        ▼
Evaluate Worker Health
        │
        ▼
Select Target
        │
        ▼
Forward Request
        │
        ▼
Receive Response
```

If the selected worker becomes unavailable, the recovery pipeline is invoked.

---

# Routing Components

| Component | Responsibility |
|-----------|----------------|
| Gateway | Entry point for requests |
| Router | Select execution target |
| Health Monitor | Evaluate worker availability |
| Recovery Manager | Coordinate failover |
| Telemetry | Record routing decisions |

Each component has a single responsibility and communicates through well-defined interfaces.

---

# Worker Selection

Worker selection is based on infrastructure state.

Current implementation considers:

- worker availability
- response status
- timeout conditions

Future implementations may additionally consider:

- latency
- resource utilization
- geographic location
- routing policies

Selection criteria should remain deterministic.

---

# Healthy Request Flow

```text
Client
   │
   ▼
Gateway
   │
   ▼
Primary Worker
   │
   ▼
Response
```

When the primary worker is healthy, no recovery action is required.

---

# Recovery Flow

```text
Client
   │
   ▼
Gateway
   │
Worker Failure
   │
   ▼
Checkpoint Lookup
   │
   ▼
Fallback Worker
   │
   ▼
Response
```

Recovery occurs without modifying application logic.

---

# Routing Decisions

Routing decisions should satisfy three conditions.

1. The selected worker is available.
2. Recovery state is preserved.
3. The routing decision is observable.

The router should not make application-specific decisions.

---

# Health Evaluation

Workers expose lightweight health information.

Typical states include:

| State | Description |
|--------|-------------|
| Healthy | Accepting requests |
| Degraded | Increased latency or intermittent failures |
| Unavailable | Unable to process requests |

Routing decisions should react only to infrastructure state.

---

# Request Metadata

The router operates on request metadata rather than application payloads.

Typical metadata includes:

| Field | Purpose |
|-------|---------|
| Request ID | Correlation |
| Agent ID | Workflow identification |
| Checkpoint ID | Recovery lookup |
| Trace ID | Distributed tracing |

Payload contents remain opaque to the routing layer.

---

# Design Constraints

The routing subsystem follows several constraints.

- Stateless operation
- Explicit routing decisions
- Framework independence
- Minimal allocations
- Observable behavior

These constraints simplify operation and improve maintainability.

---
---

# Routing Policies

The current implementation follows a deterministic routing policy.

Requests are forwarded to the primary worker whenever it is considered healthy.

```text
Request
    │
    ▼
Primary Worker
    │
Healthy?
 ┌──┴──┐
 │     │
Yes    No
 │     │
 ▼     ▼
Done  Recovery
```

Future routing policies may support:

- latency-aware routing
- weighted routing
- region-aware routing
- policy-driven routing

The routing interface is intentionally designed to allow new strategies without changing the gateway architecture.

---

# Failure Detection

The router classifies infrastructure failures before initiating recovery.

Supported recovery triggers include:

| Condition | Action |
|-----------|--------|
| HTTP 5xx | Recovery |
| Connection timeout | Recovery |
| Connection refused | Recovery |
| Worker unavailable | Recovery |

Application-specific failures are returned directly to the caller.

---

# Retry Strategy

Retries are intentionally conservative.

The routing layer does not perform unlimited retries.

Typical request flow:

```text
Request
    │
    ▼
Primary Worker
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

Repeated infrastructure failures should be surfaced through telemetry rather than hidden by continuous retries.

---

# Load Distribution

Current releases focus on correctness rather than sophisticated balancing.

The gateway routes requests to a single execution target.

Future implementations may introduce:

| Strategy | Purpose |
|----------|----------|
| Round Robin | Even request distribution |
| Least Connections | Prefer lightly loaded workers |
| Latency Aware | Minimize response time |
| Weighted Routing | Prioritize specific workers |
| Adaptive Routing | Dynamic policy selection |

These strategies are intentionally deferred until the recovery model is stable.

---

# Telemetry Integration

Every routing decision should be observable.

Typical events include:

| Event | Description |
|--------|-------------|
| Request Received | Gateway accepted a request |
| Worker Selected | Routing decision completed |
| Recovery Started | Infrastructure failure detected |
| Recovery Completed | Execution resumed |
| Request Completed | Processing finished |

Tracing and metrics should describe infrastructure behavior rather than application logic.

---

# Performance Considerations

The routing subsystem should introduce minimal overhead.

Optimization priorities include:

- efficient worker lookup
- low allocation rate
- predictable latency
- lightweight request forwarding
- inexpensive health evaluation

Benchmark-driven optimization is preferred over speculative changes.

---

# Scalability

Routing components are designed to scale independently.

```text
Clients
   │
   ▼
Load Balancer
   │
   ▼
Gateway Cluster
   │
   ▼
Worker Pool
```

Because execution state is externalized, additional gateway instances can be introduced without changing routing semantics.

---

# Operational Recommendations

Recommended operational practices include:

- monitor worker health continuously
- validate recovery paths regularly
- centralize telemetry collection
- keep routing configuration simple
- avoid unnecessary routing policies
- benchmark changes before deployment

Operational simplicity improves long-term reliability.

---

# Future Routing Work

Areas under evaluation include:

- policy-based routing
- distributed worker discovery
- adaptive health evaluation
- multi-region routing
- service mesh integration
- gRPC transport
- pluggable routing policies

These items are exploratory and are not committed implementation milestones.

---

# Design Trade-offs

The routing subsystem intentionally favors predictable behavior.

| Decision | Benefit | Cost |
|----------|---------|------|
| Stateless router | Horizontal scalability | External checkpoint dependency |
| Deterministic routing | Easier debugging | Fewer dynamic optimizations |
| Explicit recovery | Predictable behavior | Slight recovery latency |
| Small routing surface | Easier maintenance | Limited policy support |

Each trade-off prioritizes maintainability and operational clarity.

---

# Routing Summary

The routing subsystem coordinates request forwarding without executing application logic.

```text
Receive Request
        │
        ▼
Evaluate Worker
        │
        ▼
Route Request
        │
        ▼
Detect Failure
        │
        ▼
Recover If Needed
        │
        ▼
Return Response
```

The router is responsible for selecting execution targets, initiating recovery when infrastructure failures occur, and exposing routing behavior through telemetry.

Application behavior remains outside the scope of the routing layer.

---