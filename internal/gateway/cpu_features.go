// SPDX-License-Identifier: Apache-2.0
package gateway

import (
	"golang.org/x/sys/cpu"

	"github.com/devloperdevesh/FaultPlane/internal/gateway/bpf"
)

// HardwareTrack represents the dynamically steered execution path.
type HardwareTrack string

const (
	TrackPlanetScale   HardwareTrack = "Planet-Scale (AVX-512 JIT Parallel Match)"
	TrackSovereignSafe HardwareTrack = "Sovereign Fallback (Pure Go Lifecycle Loop)"
)

// HardwareProfile stores dynamic CPU register evaluation and runtime configuration.
type HardwareProfile struct {
	HasAVX512F    bool          `json:"has_avx512f"`
	HasAVX512BW   bool          `json:"has_avx512bw"`
	HasAVX512CD   bool          `json:"has_avx512cd"`
	HasAVX512DQ   bool          `json:"has_avx512dq"`
	HasAVX512VL   bool          `json:"has_avx512vl"`
	CgoAvailable  bool          `json:"cgo_available"`
	ActiveTrack   HardwareTrack `json:"active_track"`
	VectorBits    int           `json:"vector_bits"`
	ParallelLanes int           `json:"parallel_lanes"`
}

// DetectHardwareCapabilities inspects native CPU registers dynamically at boot
// before memory ring buffers or routing descriptors are loaded.
func DetectHardwareCapabilities() HardwareProfile {
	hasAVX512 := cpu.X86.HasAVX512F
	activeTrack := TrackSovereignSafe
	vectorBits := 64
	parallelLanes := 1

	if hasAVX512 && bpf.CgoEnabled && bpf.IsAVX512Compiled() {
		activeTrack = TrackPlanetScale
		vectorBits = 512
		parallelLanes = 16
	}

	return HardwareProfile{
		HasAVX512F:    cpu.X86.HasAVX512F,
		HasAVX512BW:   cpu.X86.HasAVX512BW,
		HasAVX512CD:   cpu.X86.HasAVX512CD,
		HasAVX512DQ:   cpu.X86.HasAVX512DQ,
		HasAVX512VL:   cpu.X86.HasAVX512VL,
		CgoAvailable:  bpf.CgoEnabled,
		ActiveTrack:   activeTrack,
		VectorBits:    vectorBits,
		ParallelLanes: parallelLanes,
	}
}
