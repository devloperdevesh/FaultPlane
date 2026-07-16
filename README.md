# FaultPlane

> Open-source Go reliability runtime for long-running AI agent workloads with checkpointing, failure recovery, and runtime observability.

[![Go Report Card](https://goreportcard.com)](https://goreportcard.com)
[![License](https://img.shields.io/github/license/devloperdevesh/FaultPlane)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/devloperdevesh/FaultPlane)]()
[![Issues](https://img.shields.io/github/issues/devloperdevesh/FaultPlane)](https://github.com/devloperdevesh/FaultPlane/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/devloperdevesh/FaultPlane)](https://github.com/devloperdevesh/FaultPlane/pulls)

FaultPlane is an open-source Go infrastructure runtime designed to improve the reliability of long-running AI agent workloads.

It provides execution checkpointing, failure recovery, and runtime visibility so distributed AI workloads can continue from known states instead of restarting after infrastructure failures.

---

# Why FaultPlane Exists

Modern AI agents are becoming long-running distributed systems.

A single agent workflow may involve:

- multiple tool executions
- external service dependencies
- evolving context
- long-running computation

When infrastructure failures occur, applications often lose execution progress and restart expensive workflows.

FaultPlane separates execution state from execution compute.

Instead of coupling workload progress to a single worker instance, FaultPlane maintains recovery checkpoints and provides infrastructure-level recovery mechanisms.

---

# Core Capabilities

## Checkpoint Management

FaultPlane tracks execution state required for recovery.

Capabilities:

- workflow state checkpoints
- execution metadata tracking
- recovery state management
- pluggable storage design

---

## Failure Recovery

FaultPlane detects unhealthy execution paths and redirects workloads to available workers.

Recovery flow:

```text
AI Agent Request

        |

FaultPlane Gateway

        |

Primary Worker

        |

Failure Detection

        |

Checkpoint Recovery

        |

Fallback Worker

        |

Execution Resume
```

The recovery layer operates independently from application business logic.

---

## Runtime Observability

FaultPlane provides visibility into:

- worker health
- recovery events
- execution lifecycle
- runtime metrics
- system state

---

# Architecture

```text
                  AI Agent Workload

                          |

                          |

                 FaultPlane Runtime

        -----------------------------------

        Gateway Layer

        Checkpoint Engine

        Recovery Controller

        Telemetry Pipeline

        -----------------------------------

              Worker Infrastructure
```

---

# Design Principles

## Control Plane Separation

Routing and recovery logic remain independent from workload execution.

## Framework Independence

FaultPlane is designed to work with different AI agent frameworks and execution environments.

## Observable Systems

Recovery decisions should be measurable through metrics and traces.

## Minimal Integration

Reliability improvements should not require major application rewrites.

---

# Dashboard

FaultPlane provides an operations dashboard for infrastructure visibility.

Current dashboard areas:

- Runtime Metrics
- Worker Status
- Recovery Timeline
- Checkpoint History
- Telemetry Logs
- Infrastructure Health

---

# Repository Structure

```text
FaultPlane/

├── cmd/
│   └── daemon/
│
├── internal/
│   ├── api/
│   ├── control/
│   ├── gateway/
│   ├── storage/
│   └── telemetry/
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

# Current Status

Implemented:

- Go runtime foundation
- Gateway architecture
- Control plane structure
- Checkpoint management
- Recovery workflow simulation
- Operations dashboard foundation

In Progress:

- Prometheus metrics
- OpenTelemetry integration
- Persistent checkpoint storage
- Multi-worker recovery

Planned:

- Kubernetes deployment
- Distributed checkpoint replication
- Advanced runtime instrumentation
- eBPF-based telemetry experiments

---

# Future Research Areas

FaultPlane explores advanced infrastructure research directions:

- Linux eBPF runtime telemetry
- Distributed checkpoint replication
- Adaptive recovery routing
- Low-level networking instrumentation
- Kernel-assisted observability

These areas are experimental and are not required for the core runtime.

---

# Development

## Requirements

- Go 1.23+
- Docker
- Docker Compose

Clone:

```bash
git clone https://github.com/devloperdevesh/FaultPlane.git

cd FaultPlane
```

Run environment:

```bash
docker compose up --build
```

Start runtime:

```bash
go run ./cmd/daemon
```

Run tests:

```bash
go test ./...
```

Format code:

```bash
go fmt ./...
```

Static analysis:

```bash
go vet ./...
```

---

# Contributing

Contributions are welcome.

FaultPlane is looking for contributors interested in:

- Go infrastructure
- distributed systems
- storage engines
- telemetry systems
- dashboard engineering
- deployment tooling

Before submitting a pull request:

- keep changes focused
- include tests where applicable
- update documentation
- verify builds locally

Contribution workflow:

```text
Fork Repository

        |

Create Feature Branch

        |

Implement Changes

        |

Run Tests

        |

Open Pull Request
```

Good first contribution areas:

- Prometheus metrics
- dashboard improvements
- storage backends
- telemetry integrations
- recovery testing

---

# Roadmap

## Phase 1

- Local recovery engine
- In-memory checkpoints
- Failure simulation
- Dashboard foundation

## Phase 2

- Persistent storage backend
- Distributed workers
- OpenTelemetry integration
- Prometheus metrics

## Phase 3

- Kubernetes deployment
- Multi-node recovery
- Advanced runtime instrumentation
- Production benchmarking

---

# Security

Security issues should be reported privately.

See:

`SECURITY.md`

---

# License

Licensed under the Apache License 2.0.

---

# Vision

AI agents are becoming production workloads.

FaultPlane aims to provide the reliability layer that helps AI systems recover from failures, preserve execution progress, and operate continuously at scale.
