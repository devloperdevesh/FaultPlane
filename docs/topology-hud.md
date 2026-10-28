# Regional topology and recovery HUD

Open **Topology** in the UI sidebar (`/dashboard/topology`). The page polls
`GET /api/topology`, `GET /api/workers`, and `GET /api/telemetry` every five seconds,
using `NEXT_PUBLIC_API_URL` as the API base. It is read-only: selecting nodes or
replaying the demo never starts a recovery or changes routing.

## Observations and exposure

The graph groups nodes by their optional `region` field. Set `FAULTPLANE_REGION`
on the daemon to attach an operator-provided region to its registered worker.
Nodes without that metadata remain under **Region not reported**; no region is
inferred from hostnames, IP addresses, or worker IDs. Gateway and checkpoint
nodes currently have no regional instrumentation.

Worker health is taken from reported status. CPU and memory come from the
worker endpoint. Unknown states remain unknown. The worker store now records
`worker_failed` when it observes entry into a failed state, and
`worker_recovered` on a directly observed failed-to-healthy transition. Repeated
identical observations do not emit duplicate events. Initial healthy
registration is not a recovery. These events describe observed worker health,
not proof of a kernel action or successful workflow restore.

Each event includes a UTC timestamp and string metadata:

```json
{
  "type": "worker_failed",
  "timestamp": "2026-01-01T00:00:02Z",
  "metadata": {
    "node_id": "worker-a",
    "region": "example-region",
    "status": "FAILED",
    "source": "worker_registry"
  }
}
```

The HUD also recognises explicit `connection_failed`, `tcp_connection_failed`,
`recovery_failed`, `connection_recovered`, and `recovery_success` records when
producers supply them. A `node_id` (or `worker_id`) identifies the affected node;
otherwise an explicitly supplied region identifies its reported nodes. An
unknown explicit node never falls back to an entire region. Uncorrelated
incident records remain in the timeline with unknown scope. A generic
`tcp_set_state` observation or HTTP request is not treated as a failure.

Red nodes have a failed snapshot status or a latest correlated failure record.
A correlated recovery clears that event's failure indication, while a still
failed snapshot remains failed. Later incidents are ordered by event timestamp,
not response order. The timeline shows the latest 50 recognised incident
records in the API's retained window, including timestamps and outcomes.

Arrows follow the API's directed `source -> target` edges. Amber nodes are
reachable dependencies and **potential exposure**, not measured additional
failures or evidence of causal propagation. No frame-rate, routing-latency,
recovery-time, zero-allocation, or guaranteed-recovery claims are made.

## Failure and empty states

The UI distinguishes initial loading, a successful empty snapshot, a missing
incident history, and an API/connection error. After a refresh failure it retains
the last valid graph with a **stale snapshot** warning and the last successful
fetch time. Requests time out after ten seconds; polling does not overlap.
**Retry now** cancels the current request and retries immediately. Navigation
and switching to demo abort outstanding requests and clear timers.

The scrollable graph uses native keyboard-focusable node buttons and a text
inspector. The brief dashed connection animation stops after three cycles and
is disabled for `prefers-reduced-motion: reduce`.

## Reproduce one failure-to-recovery flow

1. In `ui`, run `npm ci`, then `npm run dev`.
2. Open `/dashboard/topology`, then **Open demo replay**. The banner explicitly
   identifies all data and resource values as fictional fixtures, not live data.
3. Choose **1. Healthy**, then **2. Failure**. Demo worker A fails in
   `demo-region-a`; its checkpoint dependency is potentially exposed. Worker B
   in `demo-region-b` stays healthy. Select each node to inspect the difference.
4. Choose **3. Recovered**. Worker A is active, exposure clears, and both incident
   timestamps and outcomes remain in the timeline.
5. **Return to live telemetry** discards demo state and resumes fresh API polling.
   If the daemon is unavailable, the page shows an error rather than demo data.

`npm test` in `ui` covers those controls, API decoding, regional scope,
reordered records, cycles, unknown event scope, errors, retries, timeouts,
unmount cleanup and mode isolation. `go test ./internal/api` exercises a real
worker-store failure and recovery through both topology and telemetry HTTP
handlers; this is controlled test input, not a claimed production incident.

## Remaining instrumentation boundaries

The current daemon registers its own worker as `ACTIVE`; it does not yet perform
multi-region remote health probes. The new observation recorder emits failures
only when a real producer updates the worker store with a failing state. It does
not simulate failures on a timer. Deployments need a worker health producer or
explicit correlated event producer to observe actual failures and recovery.

Kernel rtnetlink events currently enter the kernel state enforcer separately;
they are not exposed to this HTTP telemetry store with worker/region attribution.
Workflow recovery events without node/region correlation cannot be mapped to a
regional blast radius. The store retains up to 1,000 total events and is not a
durable incident history. Multiple daemon endpoints are not aggregated here.

This implementation uses the existing HTTP APIs and an accessible SVG/HTML
graph. There is no binary WebSocket endpoint in the current API, and the HUD
does not pretend to be connected to one. Streaming, durable event correlation,
remote region discovery and large-graph WebGL rendering require subsequent
instrumentation and measured scale requirements.
