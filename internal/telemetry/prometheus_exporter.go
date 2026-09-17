package telemetry

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// EnterpriseTractionCollector exposes legacy traffic and failover metrics.
// Registry-backed runtime metrics are exported by PrometheusHandler.
type EnterpriseTractionCollector struct {
	TotalVolumeBytes uint64
	SuccessfulShunts uint64
	ActiveAgentNodes uint32
}

func NewTractionCollector() *EnterpriseTractionCollector {
	return &EnterpriseTractionCollector{}
}

func (c *EnterpriseTractionCollector) AddProcessedBytes(n uint64) {
	if c == nil {
		return
	}
	atomic.AddUint64(&c.TotalVolumeBytes, n)
}

func (c *EnterpriseTractionCollector) RecordSuccessfulShunt() {
	if c == nil {
		return
	}
	atomic.AddUint64(&c.SuccessfulShunts, 1)
}

func (c *EnterpriseTractionCollector) SetActiveAgentNodes(n uint32) {
	if c == nil {
		return
	}
	atomic.StoreUint32(&c.ActiveAgentNodes, n)
}

// ExportMetricsToPrometheus exposes legacy traffic metrics using Prometheus
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

	_, _ = fmt.Fprintf(w,
		"# HELP faultplane_processed_bytes_total Total bytes processed by the data plane.\n"+
			"# TYPE faultplane_processed_bytes_total counter\n"+
			"faultplane_processed_bytes_total %d\n",
		bytes,
	)

	_, _ = fmt.Fprintf(w,
		"# HELP faultplane_successful_failover_shunts_total Total successful failover shunts.\n"+
			"# TYPE faultplane_successful_failover_shunts_total counter\n"+
			"faultplane_successful_failover_shunts_total %d\n",
		shunts,
	)

	_, _ = fmt.Fprintf(w,
		"# HELP faultplane_active_enterprise_nodes Number of active agent nodes.\n"+
			"# TYPE faultplane_active_enterprise_nodes gauge\n"+
			"faultplane_active_enterprise_nodes %d\n",
		nodes,
	)
}

// PrometheusHandler exposes the shared Registry as Prometheus text metrics.
func PrometheusHandler(registry *Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if registry == nil {
			http.Error(w, "telemetry registry unavailable", http.StatusServiceUnavailable)
			return
		}

		kernel := registry.KernelSnapshot()
		runtimeMetrics := registry.Snapshot()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.WriteHeader(http.StatusOK)

		writeCounter := func(name, help string, value uint64) {
			_, _ = fmt.Fprintf(w,
				"# HELP %s %s\n# TYPE %s counter\n%s %d\n",
				name, help, name, name, value,
			)
		}

		writeGauge := func(name, help string, value float64) {
			_, _ = fmt.Fprintf(w,
				"# HELP %s %s\n# TYPE %s gauge\n%s %g\n",
				name, help, name, name, value,
			)
		}

		writeCounter("faultplane_kernel_events_total",
			"Total kernel events observed by FaultPlane.", kernel.Events)

		writeCounter("faultplane_kernel_transitions_total",
			"Total kernel state transitions.", kernel.Transitions)

		writeCounter("faultplane_kernel_errors_total",
			"Total kernel telemetry errors.", kernel.Errors)

		writeCounter("faultplane_kernel_enforcement_total",
			"Total kernel enforcement actions.", kernel.Enforcement)

		writeCounter("faultplane_kernel_recoveries_total",
			"Total kernel recovery actions.", kernel.Recoveries)

		writeCounter("faultplane_kernel_loader_loads_total",
			"Total kernel loader load operations.", kernel.LoaderLoads)

		writeCounter("faultplane_kernel_loader_closes_total",
			"Total kernel loader close operations.", kernel.LoaderCloses)

		writeCounter("faultplane_kernel_loader_failures_total",
			"Total kernel loader failures.", kernel.LoaderFailures)

		writeCounter("faultplane_requests_total",
			"Total runtime requests.", runtimeMetrics.Requests)

		writeCounter("faultplane_recoveries_total",
			"Total runtime recoveries.", runtimeMetrics.Recoveries)

		writeCounter("faultplane_checkpoints_total",
			"Total runtime checkpoints.", runtimeMetrics.Checkpoints)

		writeGauge("faultplane_average_latency_seconds",
			"Average runtime latency in seconds.", runtimeMetrics.AverageLatency()/1000)

		writeGauge("faultplane_cpu_usage",
			"Current runtime CPU usage.", runtimeMetrics.CPU)

		writeGauge("faultplane_memory_usage",
			"Current runtime memory usage.", runtimeMetrics.Memory)

		_, _ = fmt.Fprintf(w,
			"# HELP faultplane_kernel_last_state Current kernel enforcer state.\n"+
				"# TYPE faultplane_kernel_last_state gauge\n"+
				"faultplane_kernel_last_state{state=%q} 1\n",
			kernel.LastState,
		)
	})
}
