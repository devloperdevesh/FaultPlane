# Architectural Placement: FaultPlane vs. Erlang/OTP, CRIU, and Layer 7 Service Meshes [1.5]

This specification maps the low-level systems demarcation layer proving FaultPlane's absolute microsecond defensibility against legacy userspace paradigms .

### 1. FaultPlane (L4 Go + eBPF) vs. Erlang/Elixir (OTP)
* **Erlang/OTP Ceiling:** Operates entirely within the Layer 7 Userspace Actor Model. It relies on internal serialization matrices over the BEAM virtual machine, generating massive CPU cache thrashing and Garbage Collection latency spikes under volumetric multi-tenant AI token streams .
* **FaultPlane Monopoly:** Operates directly at the Linux Kernel Layer 4 Transport Boundary. By intercepting raw byte streams natively via eBPF sockmaps, we execute connection handovers under sub-2ms bounds without ever copying packet payloads into userspace arrays.

### 2. FaultPlane vs. CRIU (Checkpoint/Restore In Userspace)
* **CRIU Ceiling:** Enforces heavy, block-level process freezing by state-dumping active memory registers onto non-volatile disk storage (SSDs). For large-context AI workloads (e.g., vLLM or 2.8T models), this introduces seconds of block-preemption drag, destroying active TCP socket links [1.5, 1.6].
* **FaultPlane Monopoly:** Zero disk I/O overhead. We utilize a pre-allocated fixed-size circular array storing unsafe pointers natively inside `internal/storage/` to execute real-time file descriptor hot-swaps inline inside memory arrays .

### 3. FaultPlane vs. Saturated L7 Service Meshes (Envoy/Istio)
* **Service Mesh Ceiling:** Forces deep packet inspection (DPI) at the application layer, inducing heavy heap allocation overheads and token serialization choke under recursive agent reasoning loops [1.5, 1.6].
* **FaultPlane Monopoly:** Zero-code-intrusion transparent sidecar proxy. Intercepts raw Layer 4 byte arrays seamlessly, maintaining flat microsecond execution integrity .
