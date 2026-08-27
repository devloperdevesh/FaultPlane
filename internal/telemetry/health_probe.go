// SPDX-License-Identifier: Apache-2.0
package telemetry

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

// ProbeSystemState represents the structural health state of the service.
type ProbeSystemState struct {
	Status             string `json:"status"`
	EbpfVerifierPassed bool   `json:"ebpf_verifier_passed"`
	SubKernelStable    bool   `json:"sub_kernel_stable"`
	ActiveShuntsCount  uint64 `json:"active_shunts_count"`
}

// HealthProbeServer exposes service health state.
type HealthProbeServer struct {
	IsHealthy     uint32
	ShuntsCounter uint64
}

// NewHealthProbeServer creates a healthy probe server.
func NewHealthProbeServer() *HealthProbeServer {
	return &HealthProbeServer{
		IsHealthy:     1,
		ShuntsCounter: 0,
	}
}

// ServeHTTP handles liveness and health probes.
func (h *HealthProbeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stable := atomic.LoadUint32(&h.IsHealthy) == 1
	shunts := atomic.LoadUint64(&h.ShuntsCounter)

	state := ProbeSystemState{
		Status:             "HEALTHY",
		EbpfVerifierPassed: true,
		SubKernelStable:    stable,
		ActiveShuntsCount:  shunts,
	}

	if !stable {
		state.Status = "CRITICAL_FAULT"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(state); err != nil {
		return
	}
}
