# Architecture

## Overview

AgentMesh is a control plane for recovering long-running AI workloads from infrastructure failures.

Instead of coupling workflow execution to a single inference runtime, AgentMesh continuously records execution checkpoints and redirects requests to healthy workers when upstream services become unavailable.

The project separates execution recovery from application logic, allowing existing AI systems to adopt recovery without introducing framework-specific orchestration.

---

## Problem Statement

Most infrastructure assumes requests are short-lived.

AI agents are different.

A single workflow may execute for several minutes while interacting with multiple language models, vector databases, external APIs, and internal services.

When an upstream runtime crashes, becomes unreachable, or exceeds timeout limits, execution usually restarts from the beginning because intermediate state exists only in process memory.

This leads to:

- Lost execution progress
- Additional inference cost
- Longer completion time
- Reduced system reliability

AgentMesh introduces an independent recovery layer that preserves execution progress outside the application runtime.

---

# Design Goals

AgentMesh is designed around a small set of engineering principles.

| Goal | Description |
|------|-------------|
| Lightweight | Keep the control plane simple and easy to operate. |
| Recoverable | Resume execution whenever a checkpoint exists. |
| Observable | Every routing decision should be measurable. |
| Independent | Avoid coupling recovery to a specific AI framework. |
| Maintainable | Prefer explicit behavior over hidden automation. |

---

# Non Goals

The project intentionally does not provide:

- workflow orchestration
- prompt management
- model routing policies
- vector database management
- distributed scheduling
- GPU orchestration
- application business logic

These responsibilities remain outside the scope of AgentMesh.

---

# High-Level Architecture

```text
                    Client
                       │
                       │
                 HTTP / gRPC
                       │
                       ▼
              +------------------+
              |     Gateway      |
              +---------+--------+
                        │
        +---------------+---------------+
        │                               │
        ▼                               ▼
 Checkpoint Engine             Telemetry Engine
        │                               │
        +---------------+---------------+
                        │
                 Routing Decision
                        │
          +-------------+-------------+
          │                           │
          ▼                           ▼
   Primary Worker             Fallback Worker
```

The gateway does not execute workloads.

Its responsibilities are limited to:

- request routing
- checkpoint lookup
- failure detection
- recovery coordination

---

# System Components

## Gateway

The gateway is the entry point for all incoming requests.

Responsibilities include:

- forwarding requests
- monitoring upstream health
- detecting failures
- invoking recovery

The gateway remains stateless.

---

## Control Plane

The control plane coordinates recovery.

Responsibilities include:

- checkpoint management
- worker selection
- routing decisions
- recovery lifecycle

No application logic is executed here.

---

## Checkpoint Store

The checkpoint store maintains the latest successful execution state.

Current implementation:

- In-memory

Future backends may include:

- Redis
- PostgreSQL
- Distributed KV stores

---

## Telemetry

Telemetry provides operational visibility.

Supported integrations (planned):

- OpenTelemetry
- Prometheus
- Jaeger
- Grafana

Metrics influence operational decisions but do not execute application logic.

---

## Workers

Workers execute application workloads.

AgentMesh treats workers as replaceable execution targets.

Workers remain unaware of routing decisions made by the control plane.

---

# Data Flow

Normal execution:

```text
Client

 │

 ▼

Gateway

 │

 ▼

Primary Worker

 │

Checkpoint

 │

 ▼

Response
```

Failure recovery:

```text
Client

 │

 ▼

Gateway

 │

503

 │

 ▼

Checkpoint Lookup

 │

 ▼

Fallback Worker

 │

Restore Context

 │

 ▼

Continue Execution
```

Recovery occurs without restarting the workflow from its initial step.

---

# Repository Layout

```text
agentmesh/

├── cmd/
│   └── daemon/

├── internal/
│   ├── api/
│   ├── control/
│   ├── gateway/
│   ├── storage/
│   └── telemetry/

├── data-plane/
│   └── agent_sim/

├── deployments/

├── docs/

└── docker-compose.yml
```

Each package has a single responsibility.

Cross-package dependencies should remain minimal.

---

# Design Constraints

AgentMesh intentionally follows these constraints.

- Small control plane
- Stateless gateway
- Explicit recovery
- Minimal dependencies
- Observable behavior
- Framework independence
- Low operational complexity
  ---

# Request Lifecycle

Every request follows the same execution path until a recovery condition is detected.

```text
Client
   │
   ▼
Gateway
   │
   ▼
Worker Selection
   │
   ▼
Primary Worker
   │
   ▼
Checkpoint Update
   │
   ▼
Response
```

The gateway remains on the request path but does not execute application logic.

---

# Checkpoint Lifecycle

Checkpoints are updated after each successful execution step.

```text
Start Workflow
       │
       ▼
Execute Step
       │
       ▼
Persist Checkpoint
       │
       ▼
Next Step
```

A checkpoint contains only the minimum state required to continue execution.

| Field | Description |
|--------|-------------|
| Agent ID | Unique workflow identifier |
| Step | Last completed step |
| Context | Serialized execution state |
| Updated At | Timestamp of latest checkpoint |

The control plane never modifies checkpoint contents.

---

# Failure Detection

Recovery begins only after an infrastructure-level failure.

Current recovery conditions include:

- HTTP 5xx responses
- Network timeout
- Connection refused
- Worker unavailable
- Gateway timeout

Application errors such as validation failures, invalid prompts, or business logic exceptions are not handled by AgentMesh.

---

# Recovery Lifecycle

```text
Request

    │

    ▼

Primary Worker

    │

Failure

    │

    ▼

Checkpoint Lookup

    │

Worker Selection

    │

    ▼

Fallback Worker

    │

Restore Context

    │

    ▼

Resume Execution
```

Recovery uses the latest available checkpoint.

No completed execution steps are repeated.

---

# Routing Strategy

Routing decisions are based on runtime availability.

```text
Healthy Worker

        │

        ▼

Route Normally



Unhealthy Worker

        │

        ▼

Recover State

        │

        ▼

Fallback Worker
```

The routing layer does not inspect application payloads.

Its only responsibility is selecting an execution target.

---

# State Model

Execution state is treated as an immutable snapshot.

```text
Workflow

↓

Checkpoint

↓

Storage

↓

Recovery

↓

Continue
```

Only successful execution updates the checkpoint.

Failed execution never overwrites existing state.

---

# Storage Design

Current implementation:

```text
Gateway

↓

Memory Store

↓

Checkpoint
```

Future architecture:

```text
Gateway

↓

Redis

↓

Persistent Storage

↓

Recovery
```

Storage remains an implementation detail.

Recovery logic should operate independently of the storage backend.

---

# Telemetry Pipeline

Telemetry is used for operational visibility.

```text
Gateway

↓

OpenTelemetry

↓

Collector

↓

Jaeger

↓

Grafana
```

Collected signals include:

| Metric | Purpose |
|---------|----------|
| Request Count | Incoming workload |
| Recovery Count | Successful recoveries |
| Recovery Latency | Time required to resume execution |
| Active Workflows | Running agents |
| Failed Requests | Unrecoverable failures |

Telemetry does not influence application behavior.

---

# Deployment Model

Local development

```text
Gateway

↓

Primary Worker

↓

Fallback Worker

↓

Jaeger
```

Future production deployment

```text
Load Balancer

        │

        ▼

Gateway Cluster

        │

        ▼

Checkpoint Backend

        │

        ▼

Worker Pool
```

The gateway layer should remain horizontally scalable.

---

# Failure Scenarios

| Scenario | Expected Behavior |
|----------|-------------------|
| Worker timeout | Recover from latest checkpoint |
| Worker crash | Route to fallback worker |
| Temporary network failure | Retry routing decision |
| Missing checkpoint | Restart workflow |
| Gateway failure | Request fails |

The recovery mechanism only guarantees progress when a valid checkpoint exists.

---

# Performance Considerations

AgentMesh prioritizes predictable behavior over maximum throughput.

Optimization targets include:

- low allocation rate
- minimal serialization
- small gateway footprint
- predictable latency
- inexpensive checkpoint updates

Performance work should focus on P95 and P99 latency rather than average latency.

---

# Scalability

The control plane is designed to scale independently from workers.

```text
Clients

    │

    ▼

Gateway Cluster

    │

    ▼

Checkpoint Backend

    │

    ▼

Worker Pool
```

Separating routing from execution allows each layer to scale independently according to workload characteristics.

---

# Design Decisions

| Decision | Rationale |
|----------|-----------|
| Stateless gateway | Easier horizontal scaling |
| External checkpoints | Enables recovery after worker failure |
| Framework independence | Simplifies adoption |
| Separate telemetry | Keeps recovery logic isolated |
| Small control plane | Reduces operational complexity |

Every design decision favors simplicity over feature count.

---
---

# Recovery Algorithm

The recovery process follows a deterministic sequence of operations.

```text
Receive Request
        │
        ▼
Forward To Worker
        │
        ▼
Worker Healthy?
        │
   ┌────┴────┐
   │         │
 Yes         No
   │         │
   ▼         ▼
Continue   Lookup Checkpoint
   │         │
   │         ▼
   │    Select Fallback
   │         │
   │         ▼
   │   Restore Context
   │         │
   └────────►▼
        Continue Execution
```

Recovery is intentionally synchronous.

The gateway waits until a routing decision has been completed before forwarding execution.

---

# Reliability Model

AgentMesh is designed to tolerate infrastructure failures rather than application failures.

Supported recovery scenarios include:

| Failure | Recoverable |
|----------|-------------|
| Worker process crash | Yes |
| Container restart | Yes |
| Network timeout | Yes |
| HTTP 5xx | Yes |
| Connection refused | Yes |

Scenarios currently outside the recovery model:

| Failure | Recoverable |
|----------|-------------|
| Invalid application state | No |
| Corrupted checkpoint | No |
| Invalid user input | No |
| Business logic exceptions | No |

Recovery guarantees depend entirely on the availability of a valid checkpoint.

---

# Consistency Model

The current implementation follows a simple consistency model.

- A checkpoint represents the latest completed execution step.
- Checkpoints are immutable once written.
- Failed execution never replaces an existing checkpoint.
- Recovery always starts from the latest successful checkpoint.

The system favors predictable recovery behavior over aggressive optimization.

---

# Security Considerations

AgentMesh assumes that checkpoint data may contain application context.

Deployments should consider:

- encrypted transport
- authenticated gateway access
- restricted administrative endpoints
- secure checkpoint storage
- audit logging

Future releases may introduce pluggable authentication and authorization mechanisms.

---

# Operational Guidelines

Recommended deployment practices:

- Keep the gateway stateless.
- Use external storage for production checkpoints.
- Export telemetry to a centralized collector.
- Monitor recovery latency and failure rate.
- Avoid storing large payloads in checkpoints.
- Regularly validate recovery paths through failure testing.

Recovery should be tested continuously rather than assumed to work.

---

# Deployment Targets

Current target environments:

| Environment | Status |
|-------------|--------|
| Local Docker | Supported |
| Linux | Supported |
| macOS | Supported |
| Windows | Supported |

Planned environments:

| Environment | Status |
|-------------|--------|
| Kubernetes | Planned |
| Multi-node clusters | Planned |
| Managed cloud deployments | Planned |

---

# Development Principles

The repository follows a small set of engineering conventions.

## Keep Packages Focused

Each package should have a single responsibility.

Avoid creating utility packages that become shared dependency buckets.

---

## Prefer Explicit Interfaces

Interfaces should describe behavior, not implementation.

Introduce abstractions only when multiple implementations exist.

---

## Keep Dependencies Small

The gateway should remain lightweight.

External dependencies should be introduced only when they provide clear operational value.

---

## Optimize After Measurement

Performance improvements should be based on benchmark data.

Avoid speculative optimization.

---

## Prioritize Readability

Code is maintained longer than it is written.

Readable implementations are preferred over clever implementations.

---

# Future Architecture

The current implementation is intentionally minimal.

Future work may include:

```text
                    Client

                      │

                      ▼

             Gateway Cluster

                      │

          ┌───────────┴───────────┐

          ▼                       ▼

   Distributed Router      Recovery Service

          │                       │

          └───────────┬───────────┘

                      ▼

              Redis Cluster

                      │

                      ▼

             Worker Pool

                      │

                      ▼

           OpenTelemetry Collector
```

This evolution keeps the control plane independent while allowing worker infrastructure to scale horizontally.

---

# Open Questions

Areas currently under evaluation:

- persistent checkpoint replication
- distributed recovery coordination
- checkpoint compression
- adaptive routing strategies
- gRPC transport
- storage abstraction
- benchmark methodology
- failure injection framework

These topics are intentionally deferred until the current architecture stabilizes.

---

# References

The design draws inspiration from established distributed systems concepts, including:

- control planes
- service meshes
- reverse proxies
- distributed tracing
- checkpoint-based recovery
- cloud-native observability

AgentMesh adapts these ideas to long-running AI execution workloads.

---

# Document Status

| Item | Status |
|------|--------|
| Architecture | Active |
| API | Experimental |
| Recovery Model | Active Development |
| Deployment Model | Experimental |
| Storage | In-Memory |
| Telemetry | Planned |

This document reflects the current architecture and will evolve alongside the implementation.

---
