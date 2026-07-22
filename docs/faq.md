# Frequently Asked Questions

This document answers common questions about FaultPlane.

The FAQ covers architecture, deployment, recovery behavior, and operational concepts.

For deeper technical details, refer to the design documents under `docs/`.

---

# General

## What is FaultPlane?

FaultPlane is a transport resilience layer for long-running AI workloads.

It provides infrastructure-level failure detection, recovery coordination, checkpoint management, and runtime observability without requiring application code changes.

---

## What problem does FaultPlane solve?

Modern AI workloads are increasingly long-running.

An AI agent or distributed workflow may maintain:

- execution context
- external tool interactions
- streaming connections
- intermediate state
- service dependencies

for extended periods.

When infrastructure fails, execution progress is often lost because state exists only inside worker memory.

FaultPlane separates execution state from runtime processes by maintaining recoverable checkpoints.

---

## Is FaultPlane a workflow orchestration framework?

No.

FaultPlane does not orchestrate business workflows.

It provides recovery and resilience primitives.

Workflow scheduling, task management, and application logic remain responsibilities of higher-level systems.

---

## Does FaultPlane run AI models?

No.

FaultPlane does not host or execute models.

Workers remain responsible for:

- inference execution
- agent logic
- application workloads

FaultPlane manages infrastructure reliability.

---

## Is FaultPlane tied to a specific AI framework?

No.

FaultPlane is designed as an infrastructure layer independent of:

- agent frameworks
- model providers
- application languages
- workflow systems

---

# Architecture

## Why separate recovery logic from workers?

Separating recovery from execution improves fault isolation.

Workers can fail, restart, or be replaced while execution progress remains available through checkpoint storage.

---

## Why is the gateway designed to be stateless?

A stateless gateway allows multiple instances to operate simultaneously.

Recovery state is stored externally, allowing any healthy gateway instance to coordinate execution recovery.

Benefits:

- horizontal scalability
- easier deployment
- improved availability

---

## Why does FaultPlane use checkpoints?

Checkpoints preserve execution progress outside the worker lifecycle.

Without checkpoints:

```
Worker Failure

      ↓

Execution Lost

      ↓

Restart From Beginning
```

With checkpoints:

```
Worker Failure

      ↓

Restore Checkpoint

      ↓

Continue Execution
```

---

## Does FaultPlane modify application code?

No.

FaultPlane is designed as an infrastructure layer.

Applications do not require:

- custom SDKs
- framework plugins
- recovery logic
- workflow rewrites

---

# Deployment

## Can FaultPlane run locally?

Yes.

The project supports local development environments using Docker-based workflows.

Local deployments are useful for:

- development
- testing
- failure simulation

---

## Is Kubernetes required?

No.

FaultPlane can operate in simpler environments.

Kubernetes support is planned for larger production deployments requiring:

- scaling
- service discovery
- automated lifecycle management

---

## Can multiple gateway instances run together?

Yes.

The architecture supports multiple gateway instances sharing common recovery state.

Example:

```
Gateway A

      │

      ▼

Checkpoint Storage

      ▲

      │

Gateway B
```

---

# Recovery

## What failures can FaultPlane recover from?

Typical infrastructure failures include:

- worker crashes
- container restarts
- network interruptions
- request timeouts
- unhealthy runtime instances

Application-specific failures are outside the recovery boundary.

---

## What happens if no checkpoint exists?

Recovery requires valid stored state.

If no checkpoint exists, the workload must restart from its initial execution point.

---

## Does recovery change application behavior?

No.

FaultPlane restores previously recorded execution state.

Application logic remains unchanged.

---

# Storage

## Can different storage backends be used?

Yes.

The storage layer is designed around an abstraction boundary.

Future backends may include:

- Redis
- PostgreSQL
- distributed key-value stores
- custom storage implementations

---

## Why separate storage from routing?

Separating storage from routing keeps the architecture flexible.

Storage implementations can evolve without changing recovery logic.

---

# Telemetry

## What observability systems are supported?

FaultPlane is designed for OpenTelemetry-compatible observability.

Future integrations include:

- Prometheus metrics
- Jaeger tracing
- structured logging pipelines

---

## Does telemetry affect recovery decisions?

No.

Telemetry should remain observational.

Instrumentation should provide visibility without introducing execution dependencies.

---

# Performance

## What does FaultPlane optimize?

FaultPlane focuses on infrastructure reliability.

Primary optimization areas:

- recovery latency
- checkpoint performance
- routing overhead
- memory allocation
- throughput stability

---

## Does FaultPlane make AI inference faster?

No.

FaultPlane does not optimize model execution speed.

It improves reliability and continuity of long-running workloads.

---

# Security

## Does FaultPlane manage secrets?

No.

Secrets should be handled by dedicated systems such as:

- environment configuration
- secret managers
- cloud security services

---

## Should checkpoint data be protected?

Yes.

Checkpoint data may contain execution context.

Production deployments should consider:

- encryption at rest
- encrypted communication
- access controls
- audit logging

---

# Troubleshooting

## Recovery is not happening.

Check:

- worker health status
- checkpoint availability
- storage connectivity
- gateway logs
- timeout configuration

Recovery requires both failure detection and valid checkpoint state.

---

## Workloads always restart from the beginning.

Possible causes:

- checkpoints are not being created
- checkpoint storage is unavailable
- recovery state is invalid

Verify checkpoint creation and persistence.

---

## Workers are healthy but traffic fails.

Review:

- gateway configuration
- worker endpoints
- network connectivity
- telemetry events

Operational metrics usually provide the first indication of failure.

---

# Project Status

## Is FaultPlane production ready?

FaultPlane is currently under active development.

Current focus areas:

- stabilizing recovery workflows
- improving observability
- adding persistent storage
- expanding deployment support
- improving testing coverage

---

# Contributing

## How can I contribute?

See:

```
CONTRIBUTING.md
```

for:

- development workflow
- coding standards
- testing requirements
- pull request guidelines

---

## Where should architecture discussions happen?

Large design changes should be discussed before implementation.

Recommended channels:

- GitHub Discussions
- Design proposals
- Architecture issues

Early discussion helps maintain consistency.

---

# Documentation

## Where should I start?

Recommended reading order:

1. `README.md`
2. `ARCHITECTURE.md`
3. `DESIGN.md`
4. `CHECKPOINTS.md`
5. `DEPLOYMENT.md`
6. `BENCHMARKS.md`
7. `SECURITY.md`
8. `ROADMAP.md`

---

# Future Work

Planned capabilities include:

- persistent checkpoint backends
- distributed recovery
- Kubernetes deployment
- advanced telemetry
- automated benchmarking
- additional storage providers
- smarter routing strategies

Development will continue incrementally while maintaining architectural simplicity.

---

# Additional Questions

If your question is not answered here:

- review the documentation in `docs/`
- search existing GitHub discussions
- open a new discussion with relevant context

Community feedback helps improve FaultPlane's architecture and usability.
