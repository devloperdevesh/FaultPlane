# Frequently Asked Questions

This document answers common questions about AgentMesh.

The FAQ focuses on architecture, deployment, recovery, and operational behavior. For implementation details, refer to the corresponding documents in the `docs/` directory.

---

# General

## What is AgentMesh?

AgentMesh is a control plane for recovering long-running AI workloads after infrastructure failures.

It coordinates routing, checkpoint management, and recovery while remaining independent of workload execution.

---

## What problem does AgentMesh solve?

Modern AI workloads may execute for several minutes or longer.

If a worker becomes unavailable during execution, progress is often lost because execution state exists only in process memory.

AgentMesh externalizes execution progress into checkpoints so that workloads can resume from the latest successful state.

---

## Is AgentMesh an orchestration framework?

No.

AgentMesh coordinates recovery.

Workflow orchestration remains the responsibility of higher-level systems.

---

## Does AgentMesh execute AI models?

No.

Workers execute application logic.

AgentMesh coordinates infrastructure behavior.

---

## Is AgentMesh tied to a specific AI framework?

No.

The project is designed to remain independent of agent frameworks and model providers.

---

# Architecture

## Why separate the control plane from workers?

Separating recovery logic from workload execution improves fault isolation and allows each layer to evolve independently.

Workers can fail without affecting the recovery subsystem.

---

## Why is the gateway stateless?

A stateless gateway allows multiple gateway instances to share the same recovery state through external checkpoint storage.

This simplifies scaling and improves availability.

---

## Why use checkpoints?

Checkpoints preserve execution progress outside the worker process.

Recovery becomes possible even if the original worker terminates unexpectedly.

---

## Can checkpoints be stored persistently?

Yes.

The current implementation uses in-memory storage, while persistent backends such as Redis or PostgreSQL are planned.

---

# Deployment

## Can AgentMesh run locally?

Yes.

The repository includes a Docker Compose environment for local development and recovery testing.

---

## Is Kubernetes required?

No.

Local and single-host deployments remain valid development environments.

Kubernetes support is planned for production deployments.

---

## Can multiple gateways be deployed?

Yes.

The architecture is designed to support multiple stateless gateway instances sharing a common checkpoint backend.

---

# Recovery

## What failures can AgentMesh recover from?

Typical recovery scenarios include:

- worker crash
- container restart
- HTTP 5xx
- network timeout
- connection failure

Application-specific failures remain outside the recovery model.

---

## What happens if no checkpoint exists?

The workflow must restart from its initial step because there is no saved execution state to restore.

---

## Does recovery modify application data?

No.

Recovery restores previously recorded execution context.

Application behavior remains unchanged.

---

# Storage

## Can different storage backends be used?

Yes.

The storage layer is designed around an abstract interface that supports multiple backend implementations.

---

## Why separate storage from routing?

Keeping storage independent reduces coupling and allows recovery logic to remain unchanged when storage implementations evolve.

---

# Telemetry

## What telemetry standards are supported?

The project is designed around OpenTelemetry-compatible instrumentation.

Planned integrations include Prometheus and Jaeger.

---

## Does telemetry affect execution?

Telemetry is intended to remain passive.

Instrumentation should not influence routing or recovery decisions.

---

# Performance

## What should be optimized first?

The project prioritizes:

- predictable recovery latency
- low allocation rate
- efficient checkpoint operations
- stable throughput

Optimization work should be guided by benchmark results.

---

## Does AgentMesh improve model inference speed?

No.

The project focuses on infrastructure recovery rather than inference performance.

---

# Security

## Does AgentMesh manage secrets?

No.

Secrets should be managed using external secret management systems.

---

## Should checkpoint data be encrypted?

Production deployments should consider encrypting checkpoint data both in transit and at rest.

---

# Troubleshooting

## Recovery is not occurring.

Verify:

- worker health
- checkpoint availability
- storage connectivity
- timeout configuration
- gateway logs

---

## Requests always restart.

Verify that checkpoints are successfully written before failures occur.

Without valid checkpoints, recovery cannot resume execution.

---

## Workers appear healthy but routing fails.

Review:

- gateway configuration
- worker endpoints
- network connectivity
- telemetry logs

Infrastructure issues often become visible through operational metrics and traces.

---

# Project Status

## Is AgentMesh production ready?

Not yet.

The project is under active development.

Current priorities include:

- stabilizing recovery
- expanding telemetry
- introducing persistent checkpoint storage
- improving deployment support

---

# Contributing

## How can I contribute?

Refer to `CONTRIBUTING.md` for development guidelines, coding standards, testing requirements, and pull request workflow.

---

## Where should design discussions occur?

Architectural proposals are best discussed through GitHub Discussions or design issues before implementation begins.

Early discussion helps reduce implementation churn.

---

# Documentation

## Where should I start?

Recommended reading order:

1. `README.md`
2. `ARCHITECTURE.md`
3. `docs/design.md`
4. `docs/runtime.md`
5. `docs/routing.md`
6. `docs/checkpoint.md`
7. `docs/storage.md`
8. `docs/telemetry.md`
9. `docs/deployment.md`

---

# Future Work

## What major capabilities are planned?

Planned work includes:

- persistent checkpoint storage
- distributed recovery
- Kubernetes deployment
- OpenTelemetry integration
- benchmark automation
- additional storage backends
- advanced routing strategies

These items will be delivered incrementally as the implementation matures.

---

# Additional Questions

If your question is not answered here, review the documentation in the `docs/` directory or open a discussion in the project repository.