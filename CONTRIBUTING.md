# Contributing to FaultPlane

Thank you for your interest in contributing to FaultPlane.

FaultPlane is an open-source infrastructure project focused on building reliability primitives for long-running AI workloads through transport resilience, runtime recovery, and distributed system observability.

Contributions of all sizes are welcome, including:

- bug fixes
- documentation improvements
- testing improvements
- performance optimizations
- infrastructure improvements
- new runtime capabilities

This guide explains the recommended workflow and engineering standards for contributing to the project.

---

# Before You Start

Before creating an issue or pull request:

- Search existing issues and discussions.
- Review project documentation.
- Understand the current architecture.
- Confirm that the proposed change aligns with project goals.
- Keep changes focused around a single technical objective.

Large architectural changes should be discussed before implementation.

For major changes involving:

- runtime architecture
- transport layer behavior
- storage design
- recovery mechanisms
- public APIs

please open a design discussion first.

---

# Development Environment

## Requirements

| Tool | Version |
|---|---|
| Go | 1.23+ |
| Python | 3.11+ |
| Docker | Latest |
| Docker Compose | Latest |
| Git | Latest |

---

# Clone Repository

```bash
git clone https://github.com/devloperdevesh/FaultPlane.git

cd FaultPlane
```

---

# Local Development

Start supporting services:

```bash
docker compose up --build
```

Run the FaultPlane daemon:

```bash
go run ./cmd/daemon
```

Run the agent simulation environment:

```bash
python data-plane/agent_sim/main.py
```

---

# Repository Structure

```text
faultplane/

├── cmd/
│   └── daemon/
│
├── internal/
│
│   ├── api/
│   ├── gateway/
│   ├── control/
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
└── README.md
```

Each package should maintain a clear responsibility boundary.

Avoid unnecessary coupling between internal components.

---

# Development Principles

## Keep Changes Focused

A pull request should solve one primary problem.

Avoid combining:

- unrelated refactors
- feature additions
- formatting changes
- dependency updates

into a single contribution.

---

## Prefer Simple Designs

Infrastructure code should prioritize:

- correctness
- readability
- maintainability
- operational safety

Avoid introducing abstractions without a clear requirement.

---

## Preserve Architecture Boundaries

FaultPlane separates:

- data-plane execution
- recovery logic
- telemetry systems
- storage layers

Changes should respect existing package ownership.

---

## Production-Oriented Thinking

Contributions should consider:

- failure scenarios
- operational visibility
- performance impact
- backwards compatibility

---

# Coding Standards

Go code should follow standard Go practices.

Requirements:

- use `gofmt`
- handle errors explicitly
- avoid unnecessary global state
- keep interfaces small
- write clear package APIs
- prefer readable implementations

Before submitting:

```bash
go fmt ./...
```

```bash
go vet ./...
```

---

# Branch Naming

Use descriptive branch names.

Examples:

```
feature/checkpoint-storage

feature/otel-exporter

feature/runtime-metrics

fix/recovery-timeout

fix/router-state

docs/update-architecture
```

Avoid:

```
update

changes

test

final

new
```

---

# Commit Guidelines

Write commits in imperative form.

Good:

```
Add gateway health monitoring

Implement recovery state manager

Improve telemetry logging

Update architecture documentation
```

Avoid:

```
fixed stuff

changes

update

done

final commit
```

Commit history should explain the evolution of the system.

---

# Pull Request Guidelines

Every contribution should be submitted through a pull request.

Recommended workflow:

```
Create Branch

      ↓

Implement Change

      ↓

Add Tests

      ↓

Run Validation

      ↓

Update Documentation

      ↓

Open Pull Request

      ↓

Code Review

      ↓

Merge
```

---

# Pull Request Checklist

Before opening a PR:

- [ ] Code builds successfully
- [ ] Tests pass
- [ ] Formatting is applied
- [ ] Documentation is updated
- [ ] No unnecessary dependencies added
- [ ] Changes follow architecture boundaries
- [ ] Commit history is clean

---

# Testing Requirements

Run the complete test suite:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Run formatting:

```bash
go fmt ./...
```

For concurrency-sensitive changes:

```bash
go test -race ./...
```

For performance-sensitive changes:

```bash
go test -bench=. ./...
```

---

# Performance Contributions

Performance improvements should include measurements.

Useful metrics:

| Metric | Description |
|---|---|
| Throughput | Operations handled per second |
| Recovery Latency | Time required for recovery |
| Memory Usage | Runtime memory consumption |
| CPU Usage | Processing overhead |
| Allocation Count | Memory allocation behavior |

Avoid optimization without evidence.

Correctness comes before performance.

---

# Documentation Requirements

Documentation is part of the implementation.

Update documentation when changing:

- APIs
- architecture
- configuration
- deployment workflows
- runtime behavior

Technical documentation should remain:

- accurate
- concise
- implementation-focused

---

# Code Review

All changes are reviewed before merging.

Review focuses on:

- correctness
- maintainability
- architecture
- performance impact
- operational safety
- documentation quality

Review feedback should be specific and constructive.

---

# Dependency Management

Dependencies should be introduced carefully.

Before adding a dependency, consider:

- Can the standard library solve this?
- Is the dependency actively maintained?
- Does it introduce unnecessary complexity?
- Is the long-term maintenance cost acceptable?

A smaller dependency surface improves reliability and security.

---

# API Changes

Public interfaces should evolve carefully.

When changing APIs:

- document the reason
- explain migration requirements
- maintain compatibility when possible

Breaking changes should be intentional.

---

# Error Handling

Errors should provide operational context.

Prefer:

- descriptive errors
- wrapped failures
- actionable logs
- graceful recovery paths

Avoid:

- silent failures
- ignored errors
- unnecessary panics

---

# Reporting Issues

When reporting bugs, include:

- operating system
- Go version
- reproduction steps
- expected behavior
- actual behavior
- logs or error output

A minimal reproduction helps maintainers investigate quickly.

---

# Feature Requests

Feature proposals should explain:

- problem being solved
- motivation
- proposed approach
- alternatives considered
- expected impact

Large features should begin with discussion.

---

# Security Issues

Do not report security vulnerabilities through public GitHub issues.

Follow the security reporting process described in:

```
SECURITY.md
```

Responsible disclosure helps protect users and contributors.

---

# Maintainer Responsibilities

Maintainers are responsible for:

- reviewing contributions
- maintaining project quality
- managing releases
- preserving architecture consistency
- supporting contributors

Maintainers may request changes before merging.

---

# Project Scope

FaultPlane focuses on infrastructure reliability for AI-native systems.

Current scope includes:

- transport resilience
- runtime recovery
- observability
- distributed workload reliability

Out-of-scope additions may be deferred until core stability improves.

---

# Community Expectations

Contributors should:

- communicate respectfully
- provide constructive feedback
- support technical discussions
- explain decisions with evidence
- help improve the project

Engineering disagreement is healthy when handled professionally.

---

# License

By contributing to FaultPlane, you agree that your contributions will be licensed under the Apache License 2.0.

---

# Questions

For questions about contributing:

- open a GitHub discussion
- create an issue
- contact project maintainers

Thank you for helping build FaultPlane.
