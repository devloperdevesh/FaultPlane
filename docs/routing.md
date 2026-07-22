# Routing

This document describes how FaultPlane routes requests between execution runtimes.

The routing subsystem is responsible for selecting healthy execution targets, forwarding traffic, and coordinating recovery when infrastructure failures occur.

Routing does not execute workloads, interpret application behavior, or modify business state.

Its responsibility is limited to:

- request forwarding
- runtime selection
- health evaluation
- recovery coordination
- routing observability

---

# Overview

Every incoming request passes through the FaultPlane routing layer before reaching an execution runtime.

Architecture:

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
```

The routing engine makes decisions based on infrastructure state.

Application behavior remains outside the routing boundary.

---

# Design Goals

The routing subsystem is designed around the following principles.

| Goal | Description |
|---|---|
| Availability | Route traffic only to healthy runtimes. |
| Recovery | Redirect execution after infrastructure failures. |
| Predictability | Keep routing decisions deterministic. |
| Performance | Minimize routing overhead. |
| Observability | Expose routing behavior through telemetry. |

---

# Responsibilities

The routing layer is responsible for:

- selecting execution targets
- evaluating runtime health
- forwarding requests
- initiating recovery workflows
- exposing routing metrics

The routing layer is not responsible for:

- executing workloads
- workflow orchestration
- model selection
- business logic
- application state management

---

# Routing Pipeline

Every request follows a predictable routing lifecycle.

```text
Receive Request

        │

        ▼

Validate Metadata

        │

        ▼

Evaluate Runtime Health

        │

        ▼

Select Execution Target

        │

        ▼

Forward Request

        │

        ▼

Return Response
```

If the selected runtime becomes unavailable, the recovery pipeline is triggered.

---

# Routing Components

| Component | Responsibility |
|---|---|
| Gateway | External request entry point |
| Router | Execution target selection |
| Health Monitor | Runtime availability evaluation |
| Recovery Manager | Failure handling coordination |
| Telemetry Layer | Routing visibility |

Each component maintains a clear responsibility boundary.

---

# Worker Selection

Runtime selection is based on infrastructure signals.

Current considerations:

- runtime availability
- response health
- timeout conditions
- connection status

Future routing strategies may consider:

- latency
- resource utilization
- geographic placement
- workload policies

Routing decisions should remain deterministic and explainable.

---

# Healthy Request Flow

When the primary runtime is healthy:

```text
Client

   │

   ▼

Gateway

   │

   ▼

Primary Runtime

   │

   ▼

Response
```

No recovery action is required.

---

# Recovery Routing Flow

When infrastructure failure occurs:

```text
Client

   │

   ▼

Gateway

   │

   ▼

Runtime Failure Detected

   │

   ▼

Checkpoint Lookup

   │

   ▼

Recovery Runtime Selected

   │

   ▼

Execution Resumed
```

Recovery occurs without modifying application logic.

---

# Routing Decisions

Every routing decision should satisfy three conditions:

1. Selected runtime is available.
2. Execution state remains recoverable.
3. Decision is observable.

The router should never make application-specific decisions.

---

# Runtime Health Model

Workers expose lightweight health information.

Example states:

| State | Description |
|---|---|
| Healthy | Runtime can accept requests |
| Degraded | Runtime shows reduced reliability |
| Unavailable | Runtime cannot process requests |

Routing decisions should depend on infrastructure health rather than workload content.

---

# Request Metadata

The routing subsystem operates primarily on metadata.

Typical fields:

| Field | Purpose |
|---|---|
| Request ID | Request correlation |
| Workflow ID | Execution identification |
| Checkpoint ID | Recovery lookup |
| Trace ID | Distributed tracing |

Application payloads remain opaque to the routing layer.

---

# Routing Constraints

The routing subsystem follows several design constraints:

- stateless operation
- explicit decisions
- minimal latency overhead
- framework independence
- observable behavior

These constraints simplify operation and long-term maintenance.

---

# Routing Policies

The initial routing model uses deterministic selection.

Current behavior:

```text
Incoming Request

        │

        ▼

Primary Runtime

        │

        ▼

Health Check

        │

 ┌──────┴──────┐

 │             │

Healthy     Failed

 │             │

 ▼             ▼

Execute     Recovery
```

Future policies may introduce:

- round-robin routing
- weighted routing
- latency-aware routing
- region-aware routing
- adaptive policies

The routing interface is designed to support future extensions without changing the gateway architecture.

---

# Failure Detection

The router identifies infrastructure failures before starting recovery.

Supported conditions:

| Condition | Action |
|---|---|
| HTTP 5xx | Trigger recovery |
| Timeout | Trigger recovery |
| Connection failure | Trigger recovery |
| Runtime unavailable | Trigger recovery |

Application-level errors are returned without infrastructure recovery.

---

# Retry Strategy

FaultPlane avoids aggressive retry behavior.

Unlimited retries can hide infrastructure problems and increase system instability.

Typical flow:

```text
Request

   │

   ▼

Primary Runtime

   │

Failure

   │

   ▼

Recovery Attempt

   │

Success?

 ┌──┴──┐

 │     │

Yes    No

 │     │

 ▼     ▼

Done  Return Error
```

Repeated failures should become visible through telemetry.

---

# Load Distribution

Current versions prioritize correctness and recovery reliability over advanced balancing.

Future routing strategies may include:

| Strategy | Purpose |
|---|---|
| Round Robin | Even traffic distribution |
| Least Connections | Prefer less loaded runtimes |
| Latency Aware | Optimize response time |
| Weighted Routing | Control traffic distribution |
| Adaptive Routing | Dynamic decision making |

Advanced balancing will be introduced after recovery behavior is stable.

---

# Telemetry Integration

Every routing decision should produce observable signals.

Examples:

| Event | Description |
|---|---|
| Request Received | Gateway accepted traffic |
| Runtime Selected | Routing decision completed |
| Failure Detected | Runtime became unhealthy |
| Recovery Started | Failover initiated |
| Recovery Completed | Execution restored |
| Request Completed | Lifecycle finished |

Telemetry should describe infrastructure behavior.

---

# Performance Considerations

Routing should introduce minimal overhead.

Optimization priorities:

- efficient runtime lookup
- low memory allocation
- predictable latency
- lightweight health checks
- efficient request forwarding

All optimization decisions should be validated through benchmarks.

---

# Scalability

Routing components should scale independently.

Architecture:

```text
Clients

   │

   ▼

Load Balancer

   │

   ▼

FaultPlane Gateway Cluster

   │

   ▼

Worker Pool
```

Because recovery state is externalized, additional gateway instances can be added without changing routing behavior.

---

# Operational Recommendations

Recommended practices:

- monitor runtime health
- validate recovery paths
- collect routing metrics
- keep policies simple
- benchmark routing changes
- review failure patterns regularly

Operational simplicity improves reliability.

---

# Future Routing Work

Future areas include:

- distributed runtime discovery
- policy-based routing
- adaptive health evaluation
- multi-region routing
- service mesh integration
- gRPC transport
- pluggable routing engines

These remain research directions.

---

# Design Trade-offs

FaultPlane routing prioritizes predictable behavior.

| Decision | Benefit | Trade-off |
|---|---|---|
| Stateless router | Easy horizontal scaling | Requires external recovery state |
| Deterministic routing | Easier debugging | Less dynamic optimization |
| Explicit recovery | Predictable failures | Additional recovery latency |
| Minimal routing logic | Easier maintenance | Limited built-in policies |

Each decision favors reliability and operational clarity.

---

# Routing Summary

The routing subsystem coordinates traffic forwarding and recovery without executing application workloads.

```text
Receive Request

        │

        ▼

Evaluate Runtime

        │

        ▼

Select Target

        │

        ▼

Forward Request

        │

        ▼

Detect Failure

        │

        ▼

Recover If Required

        │

        ▼

Return Response
```

FaultPlane routing provides a predictable execution path while maintaining separation between:

- traffic management
- workload execution
- recovery state
- observability

This separation enables resilient long-running AI infrastructure.
