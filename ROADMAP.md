# Roadmap

This document describes the planned evolution of AgentMesh.

The roadmap is organized into incremental milestones rather than fixed release dates. Priorities may change as the project matures, implementation experience grows, and community feedback is incorporated.

AgentMesh is currently focused on establishing a reliable recovery layer for long-running AI workloads before expanding into distributed deployments.

---

# Guiding Principles

Every milestone should improve one or more of the following areas.

| Principle | Description |
|-----------|-------------|
| Reliability | Recover execution without losing completed work. |
| Simplicity | Keep the control plane small and understandable. |
| Observability | Make every recovery decision measurable. |
| Maintainability | Prefer explicit implementations over unnecessary abstractions. |
| Compatibility | Remain independent of AI frameworks whenever practical. |

---

# Current Status

| Area | Status |
|------|--------|
| Gateway | In Progress |
| Control Plane | In Progress |
| Checkpoint Engine | In Progress |
| Recovery Workflow | In Progress |
| Local Docker Environment | In Progress |
| Documentation | Active |
| Automated Testing | Planned |
| Production Deployment | Planned |

---

# Milestone 1 — Local Recovery

## Objective

Validate the recovery model on a single machine using isolated worker processes.

## Deliverables

- HTTP gateway
- Request routing
- In-memory checkpoint storage
- Local worker simulation
- Failure detection
- Automatic recovery
- Docker Compose environment
- Initial documentation

## Success Criteria

- Recovery from worker failure
- Deterministic routing behavior
- Stable local development environment

Status

**In Progress**

---

# Milestone 2 — Observability

## Objective

Expose operational visibility into the recovery process.

## Deliverables

- OpenTelemetry integration
- Prometheus metrics
- Jaeger traces
- Structured logging
- Request correlation IDs

## Metrics

- Recovery latency
- Active workflows
- Request throughput
- Worker availability
- Recovery success rate

## Success Criteria

Every routing decision should be traceable through a single execution path.

Status

**Planned**

---

# Milestone 3 — Persistent Checkpoints

## Objective

Replace in-memory checkpoint storage with durable storage backends.

## Planned Backends

| Backend | Purpose |
|----------|----------|
| Redis | Fast checkpoint lookup |
| PostgreSQL | Durable persistence |
| Pluggable Storage Interface | Backend abstraction |

## Success Criteria

Recovery remains possible after gateway restart.

Status

**Planned**

---

# Milestone 4 — gRPC Transport

## Objective

Support high-throughput communication between the gateway and workers.

## Deliverables

- gRPC gateway
- Streaming support
- Protobuf definitions
- HTTP compatibility

## Success Criteria

Reduced transport overhead while preserving recovery behavior.

Status

**Planned**

---

# Milestone 5 — Kubernetes Deployment

## Objective

Deploy AgentMesh inside Kubernetes clusters.

## Deliverables

- Kubernetes manifests
- ConfigMaps
- Secrets
- Health probes
- Rolling updates
- Horizontal scaling

## Success Criteria

Multiple gateway instances operate behind a load balancer while sharing checkpoint storage.

Status

**Planned**

---

# Milestone 6 — Distributed Recovery

## Objective

Support recovery across multiple gateway instances.

## Deliverables

- Shared checkpoint backend
- Distributed routing
- Worker discovery
- Gateway coordination

## Success Criteria

Recovery remains available even when individual gateway instances fail.

Status

**Research**

---

# Milestone 7 — Performance

## Objective

Reduce recovery overhead.

Areas of focus include:

- allocation reduction
- routing latency
- serialization overhead
- checkpoint performance
- startup time

Benchmarking will accompany every optimization.

Status

**Research**

---

# Milestone 8 — Production Readiness

The first production-ready release should include:

- persistent checkpoint storage
- OpenTelemetry support
- Prometheus metrics
- gRPC transport
- automated testing
- benchmark suite
- deployment documentation
- CI/CD pipeline
- versioned releases

Status

**Future**

---

# Testing Roadmap

Testing evolves alongside implementation.

## Unit Tests

- routing
- storage
- recovery
- gateway

---

## Integration Tests

- checkpoint restoration
- worker recovery
- timeout handling
- retry behavior

---

## Failure Injection

Simulated scenarios include:

- worker crash
- network timeout
- connection reset
- gateway restart
- checkpoint corruption

---

## Benchmarks

Measure:

| Metric | Goal |
|----------|------|
| Recovery Latency | Minimize |
| Gateway Throughput | Increase |
| Memory Allocation | Reduce |
| CPU Usage | Predictable |
| Recovery Success Rate | Maximize |

---

# Documentation Roadmap

Documentation expands with implementation.

Planned documents include:

```text
docs/

design.md

runtime.md

routing.md

checkpoint.md

telemetry.md

deployment.md

benchmarks.md

faq.md
```

Documentation should remain synchronized with the implementation.

---

# Future Research

The following areas are under evaluation.

- adaptive routing
- distributed checkpoint replication
- storage abstraction
- worker discovery
- policy engine
- eBPF-based telemetry
- recovery optimization
- multi-region deployment

Research items are exploratory and are not commitments.

---

# Out of Scope

The project does not currently plan to include:

- model hosting
- workflow orchestration
- prompt engineering
- vector databases
- model fine-tuning
- GPU scheduling
- user interfaces
- workflow builders

Keeping scope focused reduces long-term maintenance complexity.

---

# Release Strategy

Releases follow incremental milestones.

```text
Prototype

      │

      ▼

Alpha

      │

      ▼

Beta

      │

      ▼

Release Candidate

      │

      ▼

Stable
```

Breaking changes may occur before the first stable release.

---

# Long-Term Vision

AgentMesh aims to provide a lightweight recovery layer for long-running AI systems.

The project will continue prioritizing:

- predictable recovery
- operational simplicity
- observable behavior
- maintainable architecture

New functionality should strengthen these principles rather than increase feature count.

---

# Version History

| Version | Status |
|----------|--------|
| v0.x | Active Development |
| v1.0 | Planned |
