# Checkpoints

This document describes how FaultPlane captures, stores, validates, and restores execution progress.

Checkpointing is a core part of FaultPlane's recovery model.

The purpose of checkpointing is to separate execution progress from worker lifetime, allowing workloads to continue after infrastructure failures without restarting completed execution.

This document defines checkpoint semantics and recovery behavior independent of storage implementation.

---

# Overview

A checkpoint represents a durable representation of the latest valid execution state.

Instead of keeping all progress inside worker memory, FaultPlane externalizes recovery state into a checkpoint that can be restored by another runtime instance.

Architecture:

```
Execution

    │

    ▼

Checkpoint Manager

    │

    ▼

Storage Layer

    │

    ▼

Recovery Manager

    │

    ▼

Resume Execution
```

The worker executing the workload does not own the lifetime of execution state.

---

# Design Goals

The checkpoint subsystem is designed around the following principles.

| Goal | Description |
|---|---|
| Reliability | Preserve completed execution progress during failures. |
| Isolation | Separate runtime state from worker availability. |
| Consistency | Restore the latest valid checkpoint deterministically. |
| Portability | Avoid dependency on specific application frameworks. |
| Simplicity | Keep recovery semantics predictable and easy to reason about. |

---

# Scope

The checkpoint subsystem is responsible for:

- capturing execution progress
- maintaining checkpoint metadata
- storing recovery state
- validating checkpoint integrity
- restoring execution context

The checkpoint subsystem is **not responsible** for:

- running workloads
- scheduling workers
- routing network traffic
- application business logic
- user data management

---

# Checkpoint Lifecycle

A checkpoint follows a controlled lifecycle.

```
Create

  ↓

Validate

  ↓

Persist

  ↓

Reference

  ↓

Recover

  ↓

Expire
```

Only successfully completed execution steps should create new checkpoints.

Failed operations must not overwrite valid recovery state.

---

# State Model

A checkpoint represents execution state at a specific point in time.

Example metadata:

| Field | Description |
|---|---|
| Workflow ID | Unique execution identifier |
| Checkpoint ID | Unique checkpoint reference |
| Version | Schema version |
| Execution Step | Last completed step |
| Context | Serialized runtime state |
| Timestamp | Creation time |
| Metadata | Additional recovery information |

The exact schema may evolve as the recovery model matures.

---

# Ownership Model

Execution state belongs to the checkpoint subsystem, not the worker.

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

Recovery Manager

   │

   ▼

New Worker
```

A worker failure should not automatically destroy execution progress.

---

# Checkpoint Creation Policy

FaultPlane creates checkpoints only after successful execution boundaries.

Flow:

```
Execute Step

      │

      ▼

Execution Successful?

      │

 ┌────┴────┐

 │         │

Yes        No

 │         │

 ▼         ▼

Save    Ignore

Checkpoint
```

Invalid or incomplete execution states should never become recovery points.

---

# Recovery Model

Recovery starts from the latest valid checkpoint.

```
Failure Detected

        ↓

Locate Checkpoint

        ↓

Validate State

        ↓

Restore Context

        ↓

Continue Execution
```

Recovery behavior should remain deterministic.

Given the same checkpoint and runtime conditions, recovery should produce the same execution result.

---

# Storage Abstraction

Checkpoint management is independent from storage implementation.

Architecture:

```
Worker

   │

   ▼

Checkpoint Manager

   │

   ▼

Storage Interface

   │

   ▼

Storage Backend
```

Possible future storage implementations:

- in-memory store
- Redis
- PostgreSQL
- distributed key-value systems

The recovery layer should not depend on a specific database technology.

---

# Checkpoint Validation

Before restoration, checkpoints must be validated.

Validation includes:

- checkpoint identifier verification
- metadata completeness
- schema compatibility
- serialization integrity
- timestamp validation

Invalid checkpoints should never be restored.

---

# Consistency Model

FaultPlane follows a predictable recovery model.

Rules:

- latest valid checkpoint wins
- incomplete checkpoints are ignored
- successful execution is preserved
- failed execution does not overwrite state
- recovery uses immutable checkpoint data

The system prioritizes correctness over aggressive optimization.

---

# Failure Handling

FaultPlane handles common infrastructure failures.

| Failure Scenario | Expected Behavior |
|---|---|
| Worker crash | Restore latest checkpoint |
| Container restart | Restore available checkpoint |
| Network interruption | Continue from recovery state |
| Gateway restart | Recover using persistent storage |
| Missing checkpoint | Start new execution |

Recovery guarantees depend on checkpoint availability.

---

# Checkpoint Versioning

Checkpoint schemas may evolve over time.

Each checkpoint should include version information.

```
Checkpoint

      ↓

Schema Version Check

      ↓

Compatible?

 ┌────┴────┐

 │         │

Yes        No

 │         │

 ▼         ▼

Restore   Reject
```

Schema migrations should be handled explicitly.

---

# Garbage Collection

Checkpoint cleanup is required to prevent unlimited storage growth.

Possible strategies:

- workflow completion cleanup
- time-based expiration
- storage quota policies
- scheduled cleanup jobs

Cleanup operations must never remove active recovery points.

---

# Performance Considerations

Checkpoint operations should introduce minimal overhead.

Optimization priorities:

- efficient serialization
- low allocation overhead
- fast lookup
- predictable write latency
- controlled storage growth

Performance improvements should be validated through benchmarks.

---

# Security Considerations

Checkpoint data may contain execution context.

Production deployments should consider:

- encryption at rest
- encrypted transport
- access control
- audit logging
- secure deletion

Applications should avoid storing unnecessary sensitive information inside checkpoints.

---

# Operational Recommendations

Recommended practices:

- monitor checkpoint creation latency
- test recovery workflows regularly
- validate storage availability
- monitor checkpoint size growth
- maintain backup strategies
- verify recovery correctness

Reliable recovery requires both correct implementation and operational discipline.

---

# Future Improvements

Future checkpoint capabilities may include:

- persistent checkpoint backends
- incremental checkpoints
- checkpoint compression
- distributed replication
- checkpoint streaming
- pluggable storage providers
- cross-region recovery

These features will be introduced as the architecture matures.

---

# Design Trade-offs

FaultPlane intentionally favors predictable recovery behavior.

| Decision | Benefit | Trade-off |
|---|---|---|
| Immutable checkpoints | Reliable recovery | Additional storage usage |
| Externalized state | Worker independence | Storage dependency |
| Explicit validation | Safer restoration | Small validation overhead |
| Minimal metadata | Lower complexity | Reduced historical information |

These trade-offs prioritize reliability and maintainability.

---

# Checkpoint Summary

The checkpoint subsystem preserves execution progress independently from worker lifetime.

```
Execute Work

      ↓

Create Checkpoint

      ↓

Persist State

      ↓

Infrastructure Failure

      ↓

Restore Checkpoint

      ↓

Resume Execution
```

FaultPlane separates responsibilities:

- Workers execute workloads.
- Checkpoint Manager preserves progress.
- Storage maintains recovery state.
- Recovery Manager restores execution.

Together, these components provide a foundation for resilient long-running AI infrastructure.
