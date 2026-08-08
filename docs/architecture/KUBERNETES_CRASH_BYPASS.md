# Specification: Kubernetes Pod Mutation Bypass and Sub-Kernel Fault Isolation Core

This architecture specification details FaultPlane's native sub-kernel layer that completely isolates in-flight connection crashes, bypassing Kubernetes userspace API mutators and extender overheads entirely.

### 1. Zero-Allocation Line-Rate Redirection
* **Eliminating Webhook Chokes:** Unlike traditional constraints (e.g., CNCF HAMi's webhook mutations), FaultPlane operates natively within Linux Layer 4 networking structures . It tracks low-level file descriptors directly inside `internal/storage/buffer_pool.go` using lock-free Compare-And-Swap (CAS) pointers.
* **Microsecond Execution Handover:** When a distributed AI node encounters a transient network timeout, our eBPF mapas execute an instantaneous sub-2ms descriptor hot-swap, preserving volatile token context configurations natively .

### 2. Multi-Tenant Memory Fencing Alignment
* **64-Byte Cache Line Padding:** To eradicate side-channel hardware scraping risks across heterogeneous cluster tenants, all active connection variables are securely padded across discrete 64-byte boundaries, preventing L1/L2 cache line bouncing natively under peak volumetric traffic spikes .
