// Package telemetry provides runtime metrics collection
// for the AgentMesh control plane.
//
// This package is intentionally transport-agnostic.
// Metrics collected here can later be exported to:
//
//   - Prometheus
//   - OpenTelemetry
//   - Grafana
//   - Internal Dashboard
//
// No external dependency should exist inside this package.
package telemetry
