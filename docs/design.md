# Design

This document describes the architectural decisions behind FaultPlane.

Unlike the repository overview in `README.md`, this document focuses on implementation rationale, design constraints, system boundaries, and architectural trade-offs.

The goal is to explain **why FaultPlane is designed the way it is**, rather than documenting individual APIs or implementation details.

---

# Overview

FaultPlane is a transport resilience layer designed for long-running AI workloads and distributed runtime systems.

The architecture separates execution state from worker lifetime by introducing an independent recovery layer responsible for:

- failure detection
- checkpoint coordination
- recovery decisions
- runtime observability

Instead of depending on a single execution process, FaultPlane externalizes recovery state so workloads can continue after infrastructure failures.

---

# Problem Statement

Traditional distributed systems were designed around short-lived request-response interactions.

Modern AI workloads introduce different execution characteristics.

An AI workflow may maintain:

- execution context
- external tool connections
- streaming sessions
- model interactions
- intermediate state

for extended periods of time.

When infrastructure failures occur:

```
Worker Failure

      │

      ▼

Connection Lost

      │

      ▼

Execution Interrupted

      │

      ▼

Workflow Restart
```

The result is:

- repeated computation
- increased latency
- higher infrastructure cost
- reduced reliability

FaultPlane addresses this by separating workload execution from recovery management.

---

# Design Goals

FaultPlane is built around a focused set of engineering goals.

| Goal | Description |
|---|---|
| Reliability | Recover execution from the latest valid state. |
| Isolation | Separate recovery logic from workload execution. |
| Simplicity | Keep infrastructure components understandable. |
| Observability | Make recovery behavior measurable. |
| Portability | Avoid dependency on specific AI frameworks. |

---

# Non Goals

FaultPlane intentionally does not provide:

- workflow orchestration
- prompt management
- model hosting
- vector databases
- GPU scheduling
- application business logic
- AI framework abstractions

These responsibilities belong to surrounding infrastructure layers.

---

# Architectural Principles

## Separation of Concerns

Each subsystem owns a specific responsibility.

```
Gateway

→ Traffic coordination


Worker

→ Workload execution


Checkpoint Storage

→ Recovery state persistence


Telemetry

→ Operational visibility
```

Clear ownership boundaries reduce system complexity.

---

# Stateless Gateway Design

The gateway should remain stateless whenever possible.

Execution progress belongs to the checkpoint layer rather than gateway memory.

Benefits:

- horizontal scalability
- easier replacement
- simplified failure recovery
- reduced coupling

Architecture:

```
Gateway Instance A

        │

        ▼

Checkpoint Storage

        ▲

        │

Gateway Instance B
```

Any healthy gateway instance can coordinate recovery.

---

# Explicit Recovery Model

Recovery should be deterministic.

FaultPlane follows an explicit recovery sequence:

```
Failure Detection

        │

        ▼

Validate Failure

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

Continue Execution
```

Implicit retries are avoided because they make system behavior difficult to reason about.

---

# Observable Runtime

Infrastructure decisions should be measurable.

FaultPlane exposes operational visibility through:

- metrics
- traces
- structured logs
- recovery events

Important events include:

- worker failures
- routing decisions
- checkpoint operations
- recovery duration
- runtime health changes

---

# Framework Independence

FaultPlane focuses on infrastructure primitives rather than application frameworks.

The recovery layer should work with different execution environments including:

- AI agent runtimes
- inference services
- API workers
- distributed applications

Application-specific logic remains outside the system boundary.

---

# High-Level Architecture

```text
                         Client

                           │

                           ▼

                  FaultPlane Gateway

              ┌────────────┴────────────┐

              ▼                         ▼

     Failure Detection          Telemetry Pipeline

              │                         │

              └────────────┬────────────┘

                           ▼

                  Recovery Controller

                    │              │

                    ▼              ▼

             Primary Runtime   Recovery Runtime

                    │

                    ▼

             Checkpoint Storage
```

The gateway coordinates traffic.

Workers execute workloads.

Storage preserves recovery state.

---

# Control Plane

The control plane manages resilience decisions.

Responsibilities:

- runtime health evaluation
- failure detection
- checkpoint lookup
- recovery coordination
- routing decisions

The control plane does not execute application workloads.

---

# Data Plane

The data plane represents execution environments.

Examples:

- AI agent runtimes
- inference workers
- service processes
- distributed applications

Workers are replaceable execution targets.

Their lifecycle should not determine workflow durability.

---

# Recovery Workflow

Recovery follows a deterministic pipeline.

```
Incoming Request

        │

        ▼

Route To Runtime

        │

        ▼

Runtime Failure?

        │

 ┌──────┴──────┐

 │             │

 No            Yes

 │             │

 ▼             ▼

Response    Lookup Checkpoint

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

Recovery handles infrastructure failures.

Application-level failures remain outside the recovery boundary.

---

# State Ownership Model

Execution state belongs to the checkpoint subsystem.

```
Worker

   │

   ▼

Checkpoint Manager

   │

   ▼

Storage Backend

   │

   ▼

Recovery Runtime

   │

   ▼

Worker
```

Workers may fail.

Execution progress should remain recoverable.

---

# Component Boundaries

| Component | Responsibility |
|---|---|
| Gateway | Traffic coordination |
| Recovery Controller | Failure handling decisions |
| Checkpoint Manager | Execution state management |
| Storage | Persistent recovery data |
| Telemetry | Metrics and tracing |
| Worker | Workload execution |

Each component communicates through explicit interfaces.

---

# Failure Model

Current failure scenarios include:

- process crashes
- container failures
- connection interruptions
- upstream timeout
- unhealthy runtime instances

Future versions may extend recovery capabilities to additional infrastructure failures.

---

# Architectural Trade-offs

FaultPlane intentionally prioritizes reliability and maintainability.

| Decision | Benefit | Trade-off |
|---|---|---|
| Stateless gateway | Easy horizontal scaling | Requires shared state storage |
| External checkpoints | Durable recovery | Additional storage dependency |
| Framework independence | Broad compatibility | Less framework-specific optimization |
| Minimal control plane | Easier operation | Fewer built-in abstractions |

Each decision favors predictable infrastructure behavior.

---

# Evolution Strategy

FaultPlane evolves incrementally.

```
Local Recovery

        │

        ▼

Persistent Checkpoints

        │

        ▼

Distributed Recovery

        │

        ▼

Cluster Deployment

        │

        ▼

Production Infrastructure
```

Each stage introduces new capability without breaking existing architecture.

---

# Future Design Work

Areas under evaluation:

- distributed checkpoint replication
- adaptive routing
- policy-driven recovery
- gRPC transport
- Kubernetes integration
- multi-region recovery
- kernel-level observability
- recovery optimization

Research areas remain separate from committed roadmap items.

---

# Related Documentation

Additional documentation:

- `README.md` — Project overview
- `ARCHITECTURE.md` — System architecture
- `CHECKPOINTS.md` — Recovery state model
- `DEPLOYMENT.md` — Deployment patterns
- `ROADMAP.md` — Future milestones

Implementation-specific documentation is maintained under `docs/`.
