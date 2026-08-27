package telemetry

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// EnterpriseTractionCollector exposes runtime traffic and failover metrics
// in Prometheus text exposition format.
type EnterpriseTractionCollector struct {
	TotalVolumeBytes uint64
	SuccessfulShunts uint64
	ActiveAgentNodes uint32
}

// NewTractionCollector creates an empty runtime metrics collector.
func NewTractionCollector() *EnterpriseTractionCollector {
	return &EnterpriseTractionCollector{}
}

// AddProcessedBytes records bytes processed by the data plane.
func (c *EnterpriseTractionCollector) AddProcessedBytes(n uint64) {
	if c == nil {
		return
	}
	atomic.AddUint64(&c.TotalVolumeBytes, n)
}

// RecordSuccessfulShunt records a successful failover shunt.
func (c *EnterpriseTractionCollector) RecordSuccessfulShunt() {
	if c == nil {
		return
	}
	atomic.AddUint64(&c.SuccessfulShunts, 1)
}

// SetActiveAgentNodes updates the active agent-node gauge.
func (c *EnterpriseTractionCollector) SetActiveAgentNodes(n uint32) {
	if c == nil {
		return
	}
	atomic.StoreUint32(&c.ActiveAgentNodes, n)
}

// ExportMetricsToPrometheus exposes runtime metrics using the Prometheus
// text exposition format.
func (c *EnterpriseTractionCollector) ExportMetricsToPrometheus(
	w http.ResponseWriter,
	_ *http.Request,
) {
	if c == nil {
		http.Error(w, "telemetry collector unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)

	bytes := atomic.LoadUint64(&c.TotalVolumeBytes)
	shunts := atomic.LoadUint64(&c.SuccessfulShunts)
	nodes := atomic.LoadUint32(&c.ActiveAgentNodes)

	_, _ = fmt.Fprintf(
		w,
		"# HELP faultplane_processed_bytes_total Total bytes processed by the data plane.\n"+
			"# TYPE faultplane_processed_bytes_total counter\n"+
			"faultplane_processed_bytes_total %d\n",
		bytes,
	)

	_, _ = fmt.Fprintf(
		w,
		"# HELP faultplane_successful_failover_shunts_total Total successful failover shunts.\n"+
			"# TYPE faultplane_successful_failover_shunts_total counter\n"+
			"faultplane_successful_failover_shunts_total %d\n",
		shunts,
	)

	_, _ = fmt.Fprintf(
		w,
		"# HELP faultplane_active_enterprise_nodes Number of active agent nodes.\n"+
			"# TYPE faultplane_active_enterprise_nodes gauge\n"+
			"faultplane_active_enterprise_nodes %d\n",
		nodes,
	)
}
