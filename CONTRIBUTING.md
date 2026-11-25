
# Contributing to FaultPlane

Thank you for your interest in contributing to **FaultPlane**.

FaultPlane is an open-source infrastructure project focused on building reliability primitives for long-running AI workloads through transport resilience, runtime recovery, and distributed systems observability.

We welcome contributions from developers, researchers, and infrastructure engineers. Whether you are fixing a bug, improving documentation, or developing new runtime capabilities, every meaningful contribution helps strengthen the project.

This guide outlines the development workflow, engineering standards, testing expectations, and review process for contributing to FaultPlane.

---

## Table of Contents

- [Before You Start](#before-you-start)
- [Development Environment](#development-environment)
- [Getting Started](#getting-started)
- [Repository Structure](#repository-structure)
- [Development Principles](#development-principles)
- [Coding Standards](#coding-standards)
- [Branch and Commit Guidelines](#branch-and-commit-guidelines)
- [Pull Request Guidelines](#pull-request-guidelines)
- [Testing and Validation](#testing-and-validation)
- [Performance Contributions](#performance-contributions)
- [Documentation Standards](#documentation-standards)
- [Code Review](#code-review)
- [Dependency Management](#dependency-management)
- [API and Compatibility](#api-and-compatibility)
- [Error Handling](#error-handling)
- [Reporting Issues](#reporting-issues)
- [Feature Proposals](#feature-proposals)
- [Security Reporting](#security-reporting)
- [Maintainer Responsibilities](#maintainer-responsibilities)
- [Project Scope](#project-scope)
- [Community Guidelines](#community-guidelines)
- [License](#license)

---

## Before You Start

Before opening an issue or submitting a pull request:

- Search existing issues and discussions for related work.
- Review the project documentation and architecture.
- Understand the relevant components before modifying them.
- Ensure your proposed changes align with FaultPlane's objectives.
- Keep contributions focused on a clearly defined technical objective.

For significant architectural changes, open a discussion before implementation.

This particularly applies to changes involving:

- Runtime architecture
- Transport layer behavior
- Storage and buffering systems
- Recovery mechanisms
- Public APIs and interfaces

Early discussion helps avoid duplicated work and unnecessary architectural complexity.

---

## Development Environment

### Requirements

| Tool | Version |
|---|---|
| Go | 1.23+ |
| Python | 3.11+ |
| Docker | Latest stable |
| Docker Compose | Latest stable |
| Git | Latest stable |

Ensure the required tools are installed and available in your development environment.

---

## Getting Started

### 1. Fork and Clone

Fork the repository on GitHub, then clone your fork:

```bash
git clone https://github.com/<your-username>/FaultPlane.git
cd FaultPlane
```

Configure the upstream repository:

```bash
git remote add upstream https://github.com/devloperdevesh/FaultPlane.git
```

### 2. Start Supporting Services

```bash
docker compose up --build
```

### 3. Run the FaultPlane Daemon

```bash
go run ./cmd/daemon
```

### 4. Run the Agent Simulation

```bash
python data-plane/agent_sim/main.py
```

Some development commands may require additional dependencies or configuration depending on the component being worked on.

---

## Repository Structure

The following illustrates the primary project organization:

```text
FaultPlane/
│
├── cmd/
│   └── daemon/
│
├── internal/
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

Each component should maintain a clearly defined responsibility.

Contributors should avoid unnecessary coupling between internal packages and preserve existing architectural boundaries.

---

## Development Principles

### Keep Changes Focused

Each pull request should address one primary objective.

Avoid combining unrelated:

- Refactoring
- Feature implementations
- Formatting changes
- Dependency updates
- Architectural modifications

Smaller pull requests are easier to review, test, and maintain.

### Prefer Simplicity

Infrastructure code should prioritize:

- Correctness
- Readability
- Maintainability
- Predictable behavior
- Operational safety

Avoid introducing unnecessary abstractions or dependencies without a clear technical requirement.

### Preserve Architecture Boundaries

FaultPlane separates concerns across:

- Data-plane execution
- Runtime recovery
- Telemetry and observability
- Storage and buffering

Changes should respect existing package ownership and component responsibilities.

### Think About Production Behavior

Consider the following before submitting changes:

- Failure scenarios and recovery behavior
- Resource utilization
- Concurrency and synchronization
- Operational visibility
- Backward compatibility
- Security implications

---

## Coding Standards

Go contributions should follow standard Go conventions.

Requirements:

- Format Go code consistently.
- Handle errors explicitly.
- Avoid unnecessary global state.
- Keep interfaces small and focused.
- Use meaningful names.
- Write clear and maintainable implementations.
- Avoid unnecessary dependencies.

Before submitting Go changes:

```bash
go fmt ./...
go vet ./...
```

Do not introduce unrelated formatting changes into a pull request.

---

## Branch and Commit Guidelines

### Branch Naming

Use descriptive branch names that communicate the purpose of the change.

**Recommended:**

```text
feature/checkpoint-storage
feature/otel-exporter
feature/runtime-metrics

fix/recovery-timeout
fix/router-state

docs/update-architecture
```

**Avoid:**

```text
update
changes
test
final
new
```

### Commit Guidelines

Write meaningful commit messages using imperative language.

**Good examples:**

```text
Add gateway health monitoring
Implement recovery state manager
Improve telemetry logging
Update architecture documentation
```

**Avoid:**

```text
fixed stuff
changes
update
done
final commit
```

A clean commit history makes contributions easier to understand and maintain.

---

## Pull Request Guidelines

All contributions must be submitted through GitHub pull requests.

### Recommended Workflow

```text
Fork Repository
       |
Create Feature Branch
       |
Implement Changes
       |
Add or Update Tests
       |
Run Validation
       |
Update Documentation
       |
Open Pull Request
       |
Automated Checks
       |
Maintainer Review
       |
Approval and Merge
```

### Pull Request Requirements

Every pull request should:

- Clearly explain the problem being addressed.
- Describe the proposed implementation.
- Reference related issues where applicable.
- Include relevant test evidence.
- Document any behavior or API changes.
- Identify breaking changes.
- Keep modifications within the intended scope.

### Pull Request Checklist

Before requesting review:

- [ ] Changes follow project coding standards.
- [ ] Relevant tests have been added or updated.
- [ ] Applicable validation commands have been executed.
- [ ] Documentation has been updated where necessary.
- [ ] No unnecessary dependencies have been introduced.
- [ ] No unrelated files or changes are included.
- [ ] Commit history is meaningful.
- [ ] Security and compatibility implications have been considered.

### Maintainer Approval

- All pull requests require maintainer review before merging.
- Address requested changes before requesting final approval.
- Do not merge without the required repository permissions and approval.
- Maintainers may request additional tests or implementation changes.
- Contributors must not bypass required repository protections.

For fork-based contributions, maintainers will review workflow changes and permission requirements before authorizing execution.

---

## Testing and Validation

Testing requirements depend on the scope and impact of the contribution.

### Go Tests

Run the complete Go test suite:

```bash
go test ./...
```

### Static Analysis

```bash
go vet ./...
```

### Formatting

```bash
go fmt ./...
```

### Race Detection

For concurrency-sensitive changes:

```bash
go test -race ./...
```

### Benchmarks

For performance-sensitive changes:

```bash
go test -bench=. ./...
```

Run relevant targeted tests when appropriate, in addition to broader validation.

### Test Evidence

Contributors should provide verifiable validation results when requesting review.

Include:

- Commands executed
- Test results
- Relevant CI workflow status
- Logs or error output for failing checks

Do not claim that tests passed unless they were actually executed.

**Documentation-only changes:** Changes that do not modify executable code do not require unrelated Go test suites to be run. Contributors should still validate Markdown formatting, links, code examples, diagrams, and technical accuracy as applicable.

If a check cannot be executed, clearly explain why in the pull request.

---

## Performance Contributions

Performance improvements must be supported by measurable and reproducible evidence.

Relevant metrics include:

| Metric | Description |
|---|---|
| Throughput | Operations processed per second |
| Recovery Latency | Time required to restore execution |
| Memory Usage | Runtime memory consumption |
| CPU Usage | Processing overhead |
| Allocation Count | Number of memory allocations |
| Tail Latency | p95 and p99 request latency |

### Benchmark Evidence

Performance-related pull requests should include:

- Hardware and environment details
- Baseline measurements
- Optimized measurements
- Benchmark commands
- Relevant raw results
- Methodology and measurement limitations

Avoid unsupported claims about latency, throughput, or resource improvements.

**Correctness comes before performance.**

---

## Documentation Standards

Documentation is an essential part of the implementation.

Update relevant documentation when changing:

- Public APIs
- System architecture
- Configuration options
- Deployment procedures
- Runtime behavior
- Developer workflows

Technical documentation should be:

- Accurate and implementation-focused
- Clearly structured
- Consistent with existing terminology
- Free from unsupported performance claims
- Easy for contributors to understand

For mathematical documentation, ensure formulas are valid and notation is consistently defined.

For diagrams, ensure that component relationships and execution flows accurately reflect the implementation.

---

## Code Review

All contributions are subject to technical review.

Reviews focus on:

- Functional correctness
- Code quality and maintainability
- Architectural consistency
- Performance implications
- Security considerations
- Test coverage
- Documentation quality

Review feedback should be specific, constructive, and technically justified.

Contributors are expected to respond to feedback and make necessary revisions.

---

## Dependency Management

Dependencies should be introduced carefully.

Before adding a dependency, consider:

- Can the standard library solve the problem?
- Is the dependency actively maintained?
- Does it introduce security or licensing concerns?
- Does it increase operational complexity?
- Is the long-term maintenance cost justified?

Keep the dependency surface as small as reasonably possible.

---

## API and Compatibility

Public interfaces should evolve carefully.

When introducing API changes:

- Explain the motivation.
- Document behavioral differences.
- Describe migration requirements.
- Preserve backward compatibility where practical.
- Clearly identify breaking changes.

Breaking changes should be intentional and reviewed before implementation.

---

## Error Handling

Errors should provide enough context to support debugging and operational diagnosis.

Prefer:

- Descriptive error messages
- Proper error wrapping
- Actionable logs
- Explicit error handling
- Graceful recovery where appropriate

Avoid:

- Silent failures
- Ignored errors
- Unnecessary panics
- Errors without useful context

---

## Reporting Issues

When reporting a bug, provide:

- Operating system and environment
- Go or Python version, where relevant
- Reproduction steps
- Expected behavior
- Actual behavior
- Relevant logs or error output

Minimal reproducible examples are particularly helpful.

---

## Feature Proposals

Feature requests should clearly describe:

1. The problem being solved.
2. Why the capability is needed.
3. The proposed implementation.
4. Alternatives considered.
5. Expected technical impact.
6. Potential risks and compatibility implications.

Large features should begin with an architectural discussion.

---

## Security Reporting

Please do not disclose security vulnerabilities through public GitHub issues.

Follow the responsible disclosure process documented in `SECURITY.md`.

Avoid including credentials, access tokens, private keys, or other sensitive information in issues and pull requests.

Maintainers will evaluate security-related contributions with appropriate care.

---

## Maintainer Responsibilities

Maintainers are responsible for:

- Reviewing and managing contributions.
- Maintaining code quality and project standards.
- Preserving architectural consistency.
- Managing releases and project direction.
- Reviewing security-sensitive changes.
- Supporting contributors and technical discussions.

Maintainers may request modifications or decline contributions that do not meet project requirements.

---

## Project Scope

FaultPlane focuses on infrastructure reliability for AI-native systems.

Current areas of focus include:

- Transport resilience
- Runtime recovery
- Distributed workload reliability
- Observability and telemetry
- Infrastructure-level reliability primitives

Features outside the current project direction may be deferred or declined.

---

## Community Guidelines

All contributors are expected to:

- Communicate respectfully.
- Provide constructive feedback.
- Participate in technical discussions.
- Explain engineering decisions with evidence.
- Respect different technical perspectives.
- Help maintain an inclusive and collaborative environment.

Technical disagreements are welcome when handled professionally.

---

## License

By contributing to FaultPlane, you agree that your contributions will be licensed under the **Apache License 2.0**, consistent with the project's licensing terms.

---

## Questions and Support

For questions about contributing:

- Open a GitHub Discussion.
- Create an issue for project-related questions.
- Contact the project maintainers when necessary.

For bug reports and security concerns, follow their respective reporting procedures.

---

Thank you for contributing to **FaultPlane**.

Every thoughtful contribution helps build more reliable infrastructure for AI-native systems.
