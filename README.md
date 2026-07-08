# AgentMesh

> Crash recovery infrastructure for long-running AI agent workloads.

[![CI](https://img.shields.io/github/actions/workflow/status/devloperdevesh/agentmesh/ci.yml?branch=main&label=CI)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/devloperdevesh/agentmesh)]()
[![License](https://img.shields.io/github/license/devloperdevesh/agentmesh)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/devloperdevesh/agentmesh)]()
[![Issues](https://img.shields.io/github/issues/devloperdevesh/agentmesh)]()
[![Pull Requests](https://img.shields.io/github/issues-pr/devloperdevesh/agentmesh)]()

AgentMesh is a lightweight control plane for stateful AI agents. It continuously checkpoints execution state and redirects workloads to healthy workers when upstream inference nodes fail, allowing workflows to resume from the latest checkpoint instead of restarting from the beginning.

The project is designed for AI systems where execution may last minutes or hours and where restarting after transient failures is operationally expensive.

---

## Why AgentMesh Exists

Most infrastructure was designed around stateless request-response systems.

Modern AI agents are different.

An agent may execute hundreds of sequential tool calls, maintain evolving memory, and interact with multiple external systems before completing a task. If an inference server becomes unavailable halfway through execution, the application typically loses in-memory state and restarts the workflow.

AgentMesh separates execution state from execution compute.

Instead of coupling workflow progress to a single runtime instance, execution checkpoints are maintained independently so that failures can be recovered without discarding completed work.

This approach reduces unnecessary recomputation while keeping recovery logic outside the application itself.

---

## Design Goals

| Goal | Description |
|-------|-------------|
| Stateless Gateway | Keep the routing layer lightweight and independent from inference runtimes. |
| Fast Recovery | Resume execution from the latest successful checkpoint after failures. |
| Runtime Isolation | Separate control-plane logic from compute-heavy data-plane services. |
| Observable by Default | Emit traces and metrics through OpenTelemetry. |
| Language Independent | Recovery should work regardless of the agent framework being used. |
| Minimal Integration | Avoid requiring application-specific orchestration frameworks. |

---

## Non Goals

AgentMesh is intentionally not responsible for:

- workflow orchestration
- prompt engineering
- vector database management
- LLM routing policies
- distributed consensus
- model serving
- persistent memory storage
- scheduling GPU workloads

Those responsibilities belong to other systems.

---

## Architecture

```text
                   +----------------------+
                   |      AI Client       |
                   +----------+-----------+
                              |
                              |
                       HTTP / gRPC
                              |
                              v
                +---------------------------+
                |      AgentMesh Gateway    |
                +-------------+-------------+
                              |
              +---------------+----------------+
              |                                |
              |                                |
              v                                v
      Checkpoint Engine                 Telemetry Engine
              |                                |
              |                                |
              +---------------+----------------+
                              |
                     Routing Decision
                              |
                +-------------+-------------+
                |                           |
                |                           |
                v                           v
      Primary Worker               Fallback Worker
```

The gateway never owns execution.

Its responsibility is to observe runtime health, maintain checkpoint metadata, and redirect execution when failures occur.

---

## Failure Recovery

Normal execution

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
Continue
```

Failure

```text
Client
   │
   ▼
Gateway
   │
Primary Worker
   │
503 / Timeout
   │
Checkpoint Lookup
   │
Fallback Worker
   │
Resume Execution
```

Recovery happens at the infrastructure layer without requiring the application to restart from the beginning.

---

## Repository Layout

```text
agentmesh/

├── cmd/
│   └── daemon/
│
├── internal/
│   ├── api/
│   ├── control/
│   ├── gateway/
│   ├── telemetry/
│   └── storage/
│
├── data-plane/
│   └── agent_sim/
│
├── deployments/
│
├── docs/
│
├── docker-compose.yml
│
└── README.md
```

---

## Features

| Capability | Status |
|------------|--------|
| Stateful checkpointing | In Progress |
| Failure detection | In Progress |
| Automatic rerouting | In Progress |
| OpenTelemetry tracing | Planned |
| Prometheus metrics | Planned |
| Jaeger integration | Planned |
| Kubernetes deployment | Planned |
| Multi-node recovery | Planned |

---

## Quick Start

Clone the repository.

```bash
git clone https://github.com/devloperdevesh/agentmesh.git

cd agentmesh
```

Run the local environment.

```bash
docker compose up --build
```

Start the gateway.

```bash
go run ./cmd/daemon
```

Run the simulated workload.

```bash
python data-plane/agent_sim/main.py
```

The local environment starts:

| Service | Port |
|---------|------|
| Gateway | 8080 |
| Primary Worker | 8000 |
| Fallback Worker | 8001 |
| Jaeger | 16686 |

---

## Example

```

POST /route

```

```json
{
  "agent_id": "agent-001",
  "step": 14,
  "payload": {
    "context": "..."
  }
}
```

Response

```json
{
  "status": "running",
  "worker": "primary"
}
```

If the worker becomes unavailable, AgentMesh restores the latest checkpoint and redirects execution automatically.

---
---

# Architecture Principles

AgentMesh is built around a small set of architectural constraints. These constraints intentionally limit system complexity while keeping the control plane independent from application runtimes.

| Principle | Description |
|------------|-------------|
| Control Plane Separation | Routing and recovery logic execute independently from application workloads. |
| Stateless Gateway | Request routing remains lightweight. Execution state is externalized through checkpoints. |
| Failure Isolation | Worker failures must not propagate to unrelated execution streams. |
| Recovery First | Resume work whenever possible instead of restarting from the beginning. |
| Observable Systems | Recovery decisions should be backed by traces and metrics rather than assumptions. |
| Minimal Dependencies | Core runtime should remain small and easy to audit. |

---

# Core Components

| Component | Responsibility |
|-----------|----------------|
| Gateway | Entry point for all requests. Routes traffic and detects upstream failures. |
| Control Plane | Maintains checkpoint metadata and recovery state. |
| Storage | Stores execution checkpoints. Current implementation uses in-memory storage. |
| Telemetry | Publishes traces and runtime metrics. |
| Primary Worker | Normal execution target. |
| Fallback Worker | Recovery target after failures. |

---

# Execution Lifecycle

Normal execution follows four stages.

```text
Incoming Request
        │
        ▼
Gateway
        │
        ▼
Primary Worker
        │
Checkpoint Update
        │
        ▼
Continue Execution
```

If the worker becomes unavailable, recovery follows a different path.

```text
Incoming Request
        │
        ▼
Gateway
        │
503 / Timeout
        │
Checkpoint Lookup
        │
        ▼
Fallback Worker
        │
Restore Context
        │
        ▼
Resume Execution
```

Only the execution target changes.

Workflow state remains unchanged.

---

# Checkpoint Engine

The checkpoint engine records the latest successful execution state for every running workflow.

Current implementation stores checkpoints in memory.

Future versions may support persistent storage backends including Redis or distributed key-value stores.

A checkpoint contains only the information required to continue execution.

| Field | Description |
|--------|-------------|
| Agent ID | Workflow identifier |
| Step | Last completed execution step |
| Context | Serialized execution state |
| Timestamp | Last checkpoint update |

The checkpoint engine is intentionally independent from model providers and application frameworks.

---

# Failure Detection

Recovery is triggered only after explicit infrastructure failures.

Current recovery conditions include:

- HTTP 5xx responses
- connection timeout
- upstream unavailability
- network interruption

Application-level failures such as invalid prompts or business logic errors are intentionally outside the scope of AgentMesh.

---

# Routing

Routing decisions are made using runtime health information.

```text
Healthy Worker
        │
        ▼
Route Normally

Unhealthy Worker
        │
        ▼
Restore Checkpoint

        │
        ▼
Fallback Worker
```

The routing layer never modifies application state.

Its only responsibility is selecting the execution target.

---

# Observability

Every recovery decision should be observable.

The control plane exports telemetry through OpenTelemetry-compatible interfaces.

Planned integrations include:

| Integration | Purpose |
|-------------|----------|
| OpenTelemetry | Distributed tracing |
| Prometheus | Metrics collection |
| Jaeger | Trace visualization |
| Grafana | Operational dashboards |

Typical metrics include:

| Metric | Description |
|---------|-------------|
| Requests/sec | Incoming workload rate |
| Recovery Count | Successful failover events |
| Recovery Latency | Time required to restore execution |
| Active Workflows | Current running agents |
| Checkpoint Count | Stored execution checkpoints |
| Failed Requests | Requests that could not be recovered |

---

# Performance Philosophy

AgentMesh is designed to minimize recovery overhead rather than optimize inference performance.

Performance work focuses on:

- reducing recovery latency
- minimizing checkpoint overhead
- reducing gateway allocations
- avoiding unnecessary serialization
- maintaining predictable tail latency

Mean latency is less important than maintaining consistent P95 and P99 performance during failures.

---

# Trade-offs

Every system makes trade-offs.

AgentMesh intentionally chooses simplicity over feature breadth.

| Decision | Trade-off |
|-----------|-----------|
| External checkpoints | Slight checkpoint overhead in exchange for recovery capability |
| Separate control plane | Additional service, but better fault isolation |
| Lightweight gateway | Smaller feature surface, easier maintenance |
| Framework-independent design | Less framework-specific optimization |

---

# Current Status

AgentMesh is currently an early-stage infrastructure project.

Implemented:

- repository structure
- gateway skeleton
- control plane
- checkpoint management
- local recovery simulation

In progress:

- OpenTelemetry integration
- Prometheus metrics
- Jaeger tracing
- Docker deployment

Planned:

- Kubernetes deployment
- persistent checkpoint storage
- multi-node recovery
- gRPC transport
- benchmark suite

---

# Roadmap

## Phase 1

- Local recovery
- In-memory checkpoints
- Docker environment
- Failure simulation

## Phase 2

- Redis checkpoint backend
- Distributed workers
- Kubernetes deployment
- OpenTelemetry integration

## Phase 3

- Multi-region recovery
- Pluggable storage
- Adaptive routing
- eBPF-based telemetry experiments

---

# Configuration

Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| GATEWAY_PORT | 8080 | Gateway listener |
| PRIMARY_URL | localhost:8000 | Primary worker |
| FALLBACK_URL | localhost:8001 | Recovery worker |
| STORAGE | memory | Checkpoint backend |
| OTEL_EXPORTER | disabled | Telemetry exporter |

---

# Development

Run the gateway.

```bash
go run ./cmd/daemon
```

Run all tests.

```bash
go test ./...
```

Run benchmarks.

```bash
go test -bench=. ./...
```

Format source.

```bash
go fmt ./...
```

Run static analysis.

```bash
go vet ./...
```

The repository follows standard Go project layout conventions.

---
---

# Testing

AgentMesh uses standard Go tooling for testing and benchmarking.

| Task | Command |
|------|---------|
| Run unit tests | `go test ./...` |
| Run benchmarks | `go test -bench=. ./...` |
| Run race detector | `go test -race ./...` |
| Format source | `go fmt ./...` |
| Static analysis | `go vet ./...` |

As the project evolves, integration and chaos tests will be added to validate recovery behavior under simulated infrastructure failures.

---

# Project Layout

```text
agentmesh/

├── cmd/
│   └── daemon/              # Gateway entrypoint
│
├── internal/
│   ├── api/                 # HTTP handlers
│   ├── control/             # Recovery logic
│   ├── gateway/             # Routing
│   ├── storage/             # Checkpoint backend
│   └── telemetry/           # Metrics and tracing
│
├── data-plane/
│   └── agent_sim/           # Example workload
│
├── deployments/             # Deployment manifests
├── docs/                    # Design documents
├── docker-compose.yml
└── README.md
```

---

# Development

## Prerequisites

- Go 1.23+
- Python 3.11+
- Docker
- Docker Compose

Clone the repository.

```bash
git clone https://github.com/devloperdevesh/agentmesh.git

cd agentmesh
```

Start the local environment.

```bash
docker compose up --build
```

Run the gateway.

```bash
go run ./cmd/daemon
```

Run the example worker.

```bash
python data-plane/agent_sim/main.py
```

---

# Contributing

Contributions are welcome.

Before opening a pull request, please ensure that:

- the code builds successfully
- existing tests continue to pass
- new behavior includes appropriate tests when applicable
- documentation is updated if behavior changes
- commits are focused and logically grouped

## Recommended workflow

```text
Fork Repository

        │

Create Feature Branch

        │

Implement Changes

        │

Run Tests

        │

Open Pull Request
```

Use descriptive branch names.

```text
feature/checkpoint-storage
feature/grpc-gateway
fix/recovery-timeout
docs/runtime-architecture
```

---

# Pull Request Guidelines

A pull request should include:

- problem statement
- implementation summary
- testing performed
- compatibility considerations
- relevant documentation updates

Small, focused pull requests are preferred over large unrelated changes.

---

# Documentation

Additional project documentation is available under the `docs/` directory.

| Document | Description |
|-----------|-------------|
| ARCHITECTURE.md | System architecture overview |
| CONTRIBUTING.md | Development workflow |
| ROADMAP.md | Planned milestones |
| SECURITY.md | Reporting security issues |

---

# Design Philosophy

AgentMesh follows a few simple engineering principles.

- Keep the control plane small.
- Separate recovery from application logic.
- Prefer explicit behavior over hidden automation.
- Build observable systems.
- Optimize for maintainability before optimization.
- Avoid unnecessary abstractions.

The project intentionally favors clarity over feature count.

---

# Current Limitations

Current implementation intentionally keeps the scope narrow.

Known limitations include:

- in-memory checkpoint storage
- local deployment only
- HTTP transport
- single gateway instance
- simulated recovery environment

These constraints simplify development while the core recovery model matures.

---

# Future Work

Areas currently under active development include:

- persistent checkpoint storage
- gRPC support
- distributed recovery
- Kubernetes deployment
- OpenTelemetry instrumentation
- Prometheus metrics
- benchmark suite
- fault injection framework
- adaptive routing policies

---

# Frequently Asked Questions

### Is AgentMesh an orchestration framework?

No.

AgentMesh focuses on execution recovery and request routing. Workflow orchestration remains the responsibility of higher-level systems.

---

### Does AgentMesh replace Kubernetes?

No.

Kubernetes manages infrastructure resources.

AgentMesh manages execution recovery.

The two systems solve different problems and are intended to complement each other.

---

### Does AgentMesh depend on a specific AI framework?

No.

The project is designed to remain framework independent.

Any system capable of exposing execution state can integrate with the recovery layer.

---

### Is persistent storage required?

Not today.

The current implementation uses in-memory checkpoints.

Persistent storage is planned as an optional backend.

---

# Versioning

Until a stable release is published, breaking changes may occur between minor versions.

Semantic Versioning will be adopted beginning with the first stable release.

---

# License

Licensed under the Apache License 2.0.

See the `LICENSE` file for details.

---

# Acknowledgements

AgentMesh is inspired by operational practices from distributed systems, service meshes, and cloud-native infrastructure.

The project builds upon established ideas in observability, fault tolerance, and runtime isolation while adapting them for long-running AI workloads.
