// SPDX-License-Identifier: Apache-2.0
package gateway

import (
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/gateway/bpf"
)

var (
	// ErrRouteNotFound indicates no matching IP address was identified in the routing matrix.
	ErrRouteNotFound = errors.New("faultplane [gateway]: target routing IP not matched in allocation pool")

	// ErrInvalidRoutePool indicates an empty or uninitialized routing allocation pool.
	ErrInvalidRoutePool = errors.New("faultplane [gateway]: routing allocation pool is empty or invalid")
)

// RouteEntry represents a Layer 4 routing allocation descriptor.
type RouteEntry struct {
	IP        uint32 `json:"ip"`
	IPString  string `json:"ip_string"`
	TargetURL string `json:"target_url"`
	Healthy   bool   `json:"healthy"`
}

// RoutingEngine coordinates the multi-tiered conditional native execution framework
// across Planet-Scale (AVX-512) and Sovereign Fallback (Pure Go) tracks.
type RoutingEngine struct {
	profile       HardwareProfile
	routes        atomic.Pointer[[]RouteEntry]
	ipList        atomic.Pointer[[]uint32]
	activeTrack   HardwareTrack
	swapsTotal    atomic.Uint64
	evalLatencyNs atomic.Uint64
}

// NewRoutingEngine initializes the routing engine with boot-time dynamic hardware feature detection.
func NewRoutingEngine() *RoutingEngine {
	profile := DetectHardwareCapabilities()
	engine := &RoutingEngine{
		profile:     profile,
		activeTrack: profile.ActiveTrack,
	}

	emptyRoutes := make([]RouteEntry, 0)
	emptyIPs := make([]uint32, 0)
	engine.routes.Store(&emptyRoutes)
	engine.ipList.Store(&emptyIPs)

	return engine
}

// Profile returns the evaluated hardware registry profile.
func (e *RoutingEngine) Profile() HardwareProfile {
	return e.profile
}

// ActiveTrack returns the current execution track name.
func (e *RoutingEngine) ActiveTrack() HardwareTrack {
	return e.activeTrack
}

// TotalEvaluations returns the total routing evaluations processed.
func (e *RoutingEngine) TotalEvaluations() uint64 {
	return e.swapsTotal.Load()
}

// LastLatencyNs returns the duration in nanoseconds of the most recent evaluation.
func (e *RoutingEngine) LastLatencyNs() uint64 {
	return e.evalLatencyNs.Load()
}

// SetRoutes atomically replaces the active routing table descriptors without mutex contention.
func (e *RoutingEngine) SetRoutes(entries []RouteEntry) {
	newRoutes := make([]RouteEntry, len(entries))
	copy(newRoutes, entries)

	newIPs := make([]uint32, len(entries))
	for i, r := range entries {
		newIPs[i] = r.IP
		if newRoutes[i].IPString == "" {
			newRoutes[i].IPString = Uint32ToIP(r.IP).String()
		}
	}

	e.routes.Store(&newRoutes)
	e.ipList.Store(&newIPs)
}

// MatchIP executes high-throughput IP matching against active routing allocations.
// In Track 1 (Planet-Scale), it cross-evaluates 16 IP addresses in parallel via AVX-512.
// In Track 2 (Sovereign Fallback), it executes contiguous sequential scanning loops in pure Go.
func (e *RoutingEngine) MatchIP(targetIP uint32) (RouteEntry, error) {
	start := time.Now()
	defer func() {
		dur := uint64(time.Since(start).Nanoseconds())
		e.evalLatencyNs.Store(dur)
		e.swapsTotal.Add(1)
	}()

	ipsPtr := e.ipList.Load()
	routesPtr := e.routes.Load()
	if ipsPtr == nil || routesPtr == nil || len(*ipsPtr) == 0 {
		return RouteEntry{}, ErrInvalidRoutePool
	}

	ips := *ipsPtr
	routes := *routesPtr

	matchIdx := -1

	// Track 1: Planet-Scale Track (AVX-512 JIT Parallel Match)
	if e.activeTrack == TrackPlanetScale && len(ips) >= 16 {
		matchIdx = bpf.NativeEvaluateIPs(ips, targetIP)
	}

	// Track 2: Sovereign Fallback Track (Pure Go Lifecycle Loop)
	// Guaranteed execution path when SIMD extensions or Cgo are unavailable, or for sub-16 tails.
	if matchIdx < 0 {
		matchIdx = e.evaluateSequentialFallback(ips, targetIP)
	}

	if matchIdx >= 0 && matchIdx < len(routes) {
		return routes[matchIdx], nil
	}

	return RouteEntry{}, ErrRouteNotFound
}

// evaluateSequentialFallback provides the air-gapped safe fallback sequential scanning loop.
func (e *RoutingEngine) evaluateSequentialFallback(ips []uint32, targetIP uint32) int {
	n := len(ips)
	for i := 0; i < n; i++ {
		if ips[i] == targetIP {
			return i
		}
	}
	return -1
}

// IPToUint32 converts a standard IPv4 net.IP into a big-endian uint32.
func IPToUint32(ip net.IP) uint32 {
	ip4 := ip.To4()
	if ip4 == nil {
		return 0
	}
	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
}

// Uint32ToIP converts a big-endian uint32 back to net.IP.
func Uint32ToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}
