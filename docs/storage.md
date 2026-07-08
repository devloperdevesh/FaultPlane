# Storage

This document describes the storage subsystem used by AgentMesh.

The storage layer persists execution checkpoints and exposes a simple interface to the recovery system. It is intentionally isolated from routing, runtime execution, and telemetry.

Storage exists to preserve workflow progress across infrastructure failures while remaining independent of worker processes.

---

# Overview

The storage subsystem is responsible for persisting checkpoint state.

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
Storage Backend
```

The storage implementation is replaceable.

Recovery logic should not depend on a specific backend.

---

# Design Goals

The storage layer is designed around the following principles.

| Goal | Description |
|------|-------------|
| Durability | Preserve checkpoint state reliably. |
| Simplicity | Keep the storage interface minimal. |
| Portability | Support multiple storage implementations. |
| Predictability | Provide deterministic read and write behavior. |
| Independence | Separate persistence from recovery logic. |

---

# Scope

The storage subsystem is responsible for:

- writing checkpoints
- reading checkpoints
- deleting expired checkpoints
- validating stored metadata

The storage subsystem is not responsible for:

- request routing
- worker execution
- telemetry collection
- workflow scheduling
- business logic

---

# Storage Architecture

```text
Gateway
    │
    ▼
Checkpoint Manager
    │
    ▼
Storage Interface
    │
    ▼
Backend
```

The interface between the checkpoint manager and the backend should remain stable even if storage implementations change.

---

# Storage Interface

The storage interface should expose a small set of operations.

| Operation | Purpose |
|-----------|---------|
| Save | Persist a checkpoint |
| Load | Retrieve a checkpoint |
| Delete | Remove a checkpoint |
| Exists | Verify checkpoint presence |
| List | Enumerate checkpoints (optional) |

Additional operations should only be introduced when required by the recovery model.

---

# Backend Abstraction

Storage implementations should satisfy the same interface.

Possible backends include:

| Backend | Status |
|----------|--------|
| In-Memory | Current |
| Redis | Planned |
| PostgreSQL | Planned |
| Distributed KV Store | Research |

The recovery pipeline should not require changes when switching storage backends.

---

# Data Model

A checkpoint record contains the minimum information required to resume execution.

| Field | Description |
|--------|-------------|
| Agent ID | Workflow identifier |
| Checkpoint ID | Unique record identifier |
| Step | Last completed execution step |
| Context | Serialized execution state |
| Timestamp | Creation time |

Storage implementations may include additional metadata for operational purposes.

---

# Read Operations

The storage layer returns the latest valid checkpoint for a workflow.

```text
Agent ID
    │
    ▼
Lookup
    │
    ▼
Checkpoint
```

Read operations should be deterministic and low latency.

---

# Write Operations

Checkpoint updates occur after successful execution.

```text
Execute Step
      │
      ▼
Serialize
      │
      ▼
Persist
      │
      ▼
Return
```

Failed writes should never replace an existing valid checkpoint.

---

# Delete Operations

Checkpoints may be removed after:

- workflow completion
- expiration
- administrative cleanup

Deletion policies should avoid removing checkpoints required for active workflows.

---

# Consistency

The storage subsystem follows a simple consistency model.

- latest valid checkpoint is authoritative
- completed work is preserved
- failed writes are discarded
- checkpoint history is append-oriented where practical

Deterministic behavior is preferred over aggressive optimization.

---

# Design Constraints

The storage subsystem follows several constraints.

- backend independent
- deterministic reads
- explicit writes
- minimal metadata
- predictable latency

These constraints simplify maintenance and future backend expansion.

---
---

# Persistence Strategy

The storage subsystem separates persistence from execution.

Workers never interact with storage directly.

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
Storage Backend
```

This separation allows storage implementations to evolve independently of the execution pipeline.

---

# Storage Lifecycle

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
Delete
```

Each stage should be deterministic and independently testable.

---

# Failure Handling

Storage failures should be handled explicitly.

| Failure | Expected Behavior |
|----------|-------------------|
| Backend unavailable | Return error to recovery manager |
| Read timeout | Retry according to recovery policy |
| Write timeout | Preserve previous checkpoint |
| Corrupted record | Reject checkpoint |
| Missing checkpoint | Restart workflow |

The storage subsystem should never silently discard errors.

---

# Validation

Checkpoint metadata should be validated before persistence and before restoration.

Typical validation includes:

- identifier format
- checkpoint version
- timestamp
- serialized payload
- required metadata

Validation failures should prevent the checkpoint from entering the recovery pipeline.

---

# Versioning

Storage formats evolve independently from runtime behavior.

```text
Checkpoint
      │
      ▼
Read Version
      │
      ▼
Supported?
 ┌────┴────┐
 │         │
Yes        No
 │         │
 ▼         ▼
Load    Reject
```

Explicit versioning enables backward-compatible migrations when storage schemas change.

---

# Garbage Collection

Storage implementations should remove obsolete checkpoints in a controlled manner.

Possible strategies include:

- workflow completion
- age-based expiration
- scheduled cleanup
- storage quota enforcement

Cleanup should never remove checkpoints that may still be required for recovery.

---

# Security Considerations

Checkpoint storage may contain execution context.

Production deployments should consider:

- encryption at rest
- encrypted transport
- access control
- audit logging
- secure deletion

Applications should avoid storing secrets or unnecessary sensitive data inside checkpoint payloads.

---

# Performance Considerations

Storage performance directly affects recovery latency.

Optimization priorities include:

- efficient serialization
- low write latency
- predictable read performance
- compact checkpoint representation
- minimal allocation overhead

Performance improvements should be supported by benchmark results.

---

# Operational Recommendations

Recommended operational practices include:

- monitor storage latency
- validate backup procedures
- verify checkpoint integrity
- monitor storage capacity
- test recovery workflows regularly

Operational reliability depends on both implementation quality and deployment practices.

---

# Future Storage Work

Planned improvements include:

- persistent checkpoint storage
- Redis backend
- PostgreSQL backend
- pluggable storage providers
- checkpoint compression
- distributed replication
- incremental checkpoint support

These enhancements will be introduced incrementally as the control plane matures.

---

# Design Trade-offs

The storage subsystem intentionally favors correctness over complexity.

| Decision | Benefit | Cost |
|----------|---------|------|
| Minimal interface | Easier maintenance | Fewer backend-specific features |
| Backend abstraction | Storage flexibility | Additional interface layer |
| Explicit validation | Safer recovery | Small validation overhead |
| Deterministic reads | Predictable behavior | Less implementation freedom |

These trade-offs prioritize long-term maintainability.

---

# Storage Summary

The storage subsystem provides durable persistence for execution checkpoints while remaining independent of routing and execution.

```text
Checkpoint
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
```

Workers generate execution state.

The checkpoint manager coordinates persistence.

The storage layer provides durable state required for recovery.

Together, these components ensure execution progress survives infrastructure failures.

---