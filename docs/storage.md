# Storage

This document describes the storage subsystem used by FaultPlane.

The storage layer is responsible for preserving execution checkpoints required for recovery.

It is intentionally isolated from request routing, workload execution, and telemetry systems.

The primary goal of the storage subsystem is to ensure execution progress survives infrastructure failures without coupling recovery behavior to a specific storage backend.

---

# Overview

FaultPlane separates execution state from worker lifetime.

Workers generate checkpoint data.

The checkpoint manager coordinates persistence.

The storage layer provides durable state required during recovery.

```text
Worker Runtime

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

The storage implementation is replaceable.

Recovery logic should depend only on the storage contract, not the underlying backend.

---

# Design Goals

The storage subsystem follows several core principles.

| Goal | Description |
|---|---|
| Reliability | Preserve valid execution state during failures. |
| Simplicity | Maintain a small and predictable storage interface. |
| Portability | Support multiple backend implementations. |
| Consistency | Provide deterministic checkpoint reads and writes. |
| Isolation | Keep persistence independent from runtime execution. |

---

# Scope

The storage subsystem is responsible for:

- storing checkpoints
- retrieving checkpoints
- validating checkpoint metadata
- removing obsolete checkpoints
- maintaining recovery state

The storage subsystem is not responsible for:

- request routing
- worker execution
- workload scheduling
- telemetry collection
- application business logic

---

# Storage Architecture

The storage layer sits between execution recovery and physical persistence.

```text
                 Recovery System

                       │

                       ▼

              Checkpoint Manager

                       │

                       ▼

              Storage Interface

                       │

          ┌────────────┼────────────┐

          ▼            ▼            ▼

       Memory       Redis      PostgreSQL
```

The recovery system communicates only through the storage interface.

Backend changes should not require modifications to recovery logic.

---

# Storage Interface

The storage interface intentionally exposes minimal operations.

| Operation | Purpose |
|---|---|
| Save | Persist checkpoint state |
| Load | Retrieve checkpoint state |
| Delete | Remove obsolete checkpoints |
| Exists | Verify checkpoint availability |
| List | Inspect stored checkpoints |

Additional operations should only be introduced when required by recovery semantics.

---

# Backend Implementations

FaultPlane supports interchangeable storage implementations.

| Backend | Status |
|---|---|
| In-Memory Storage | Current |
| Redis | Planned |
| PostgreSQL | Planned |
| Distributed KV Storage | Research |

The recovery model remains unchanged regardless of backend selection.

---

# Checkpoint Data Model

A checkpoint contains the minimum information required to resume execution.

| Field | Description |
|---|---|
| Workflow ID | Identifies execution flow |
| Checkpoint ID | Unique checkpoint reference |
| Execution Step | Latest completed operation |
| Context | Serialized execution state |
| Created At | Checkpoint creation timestamp |
| Version | Schema compatibility identifier |

Additional metadata may be added for operational requirements.

---

# Checkpoint Lifecycle

Every checkpoint follows a predictable lifecycle.

```text
Create

  │

  ▼

Validate

  │

  ▼

Persist

  │

  ▼

Retrieve

  │

  ▼

Recover

  │

  ▼

Expire
```

Only successful execution states should become recovery candidates.

---

# Write Path

Checkpoint writes occur after successful execution progress.

```text
Execute Operation

        │

        ▼

Generate Checkpoint

        │

        ▼

Validate State

        │

        ▼

Persist

        │

        ▼

Confirm Success
```

Failed writes must not overwrite the last valid checkpoint.

---

# Read Path

Recovery begins by locating the latest valid checkpoint.

```text
Workflow ID

      │

      ▼

Checkpoint Lookup

      │

      ▼

Validate Metadata

      │

      ▼

Restore State
```

Reads should be deterministic and optimized for recovery latency.

---

# Consistency Model

The storage subsystem follows a simple consistency model.

Principles:

- latest valid checkpoint is authoritative
- invalid checkpoints are rejected
- failed writes preserve previous state
- recovery uses validated state only

Correctness is prioritized over aggressive optimization.

---

# Persistence Strategy

Storage is separated from execution.

Workers never directly manage persistence.

```text
Worker

   │

   ▼

Checkpoint Manager

   │

   ▼

Storage Interface

   │

   ▼

Persistent Backend
```

This separation allows storage technology to evolve independently.

---

# Failure Handling

Storage failures are explicitly handled.

| Failure | Expected Behavior |
|---|---|
| Backend unavailable | Recovery request fails safely |
| Read timeout | Retry according to policy |
| Write failure | Preserve previous checkpoint |
| Corrupted checkpoint | Reject restoration |
| Missing checkpoint | Restart execution |

Storage failures should never be silently ignored.

---

# Checkpoint Validation

Before storing or restoring checkpoints, metadata validation should occur.

Validation includes:

- identifier correctness
- schema version compatibility
- timestamp validity
- serialized data integrity
- required metadata presence

Invalid checkpoints must never enter the recovery pipeline.

---

# Versioning

Checkpoint schemas evolve over time.

Version identifiers allow controlled migrations.

```text
Checkpoint

      │

      ▼

Read Version

      │

      ▼

Compatible?

 ┌────┴────┐

 │         │

Yes        No

 │         │

 ▼         ▼

Load    Reject
```

Schema evolution should remain backward compatible whenever practical.

---

# Garbage Collection

Obsolete checkpoints should be removed safely.

Possible cleanup strategies:

- completed workflow cleanup
- age-based expiration
- storage quota management
- scheduled garbage collection

Active recovery checkpoints must never be removed.

---

# Security Considerations

Checkpoint data may contain execution context.

Production deployments should consider:

- encryption at rest
- encrypted transport
- access control
- audit logging
- secure deletion

Applications should avoid storing unnecessary sensitive information.

---

# Performance Considerations

Storage performance directly impacts recovery speed.

Optimization priorities:

- low checkpoint write latency
- efficient serialization
- compact state representation
- predictable read performance
- reduced memory allocation

All optimizations should be validated through benchmarks.

---

# Operational Recommendations

Recommended practices:

- monitor storage latency
- validate checkpoint integrity
- test recovery workflows
- monitor storage growth
- verify backup procedures
- review retention policies

Reliable recovery requires both correct implementation and operational discipline.

---

# Future Storage Work

Future improvements include:

- persistent production backends
- Redis integration
- PostgreSQL integration
- distributed checkpoint replication
- checkpoint compression
- incremental checkpoints
- storage provider plugins

These capabilities will be introduced as the recovery architecture matures.

---

# Design Trade-offs

The storage design intentionally favors reliability and simplicity.

| Decision | Benefit | Cost |
|---|---|---|
| Minimal interface | Easier maintenance | Limited backend-specific features |
| Backend abstraction | Storage flexibility | Additional interface layer |
| Explicit validation | Safer recovery | Small validation overhead |
| Deterministic reads | Predictable recovery | Less optimization freedom |

These trade-offs improve long-term maintainability.

---

# Storage Summary

The storage subsystem preserves execution progress independently from worker lifetime.

```text
Execution State

        │

        ▼

Checkpoint

        │

        ▼

Validate

        │

        ▼

Persist

        │

        ▼

Recover
```

Workers produce execution state.

The checkpoint manager coordinates persistence.

The storage layer maintains recovery state.

Together, these components allow FaultPlane to recover long-running AI workloads without losing completed execution progress.
