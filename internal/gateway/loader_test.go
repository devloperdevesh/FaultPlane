// SPDX-License-Identifier: Apache-2.0
package gateway

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestRoutingEngine_EmptyPool(t *testing.T) {
	engine := NewRoutingEngine()
	_, err := engine.MatchIP(12345)
	if err != ErrInvalidRoutePool {
		t.Fatalf("expected ErrInvalidRoutePool, got %v", err)
	}
}

func TestRoutingEngine_FallbackSequential(t *testing.T) {
	engine := NewRoutingEngine()

	// Enforce fallback track for explicit test verification
	engine.activeTrack = TrackSovereignSafe

	routes := make([]RouteEntry, 32)
	for i := 0; i < 32; i++ {
		ip := uint32(10<<24 | 0<<16 | 0<<8 | (i + 1))
		routes[i] = RouteEntry{
			IP:        ip,
			TargetURL: fmt.Sprintf("http://worker-%d.internal:8080", i+1),
			Healthy:   true,
		}
	}
	engine.SetRoutes(routes)

	// Test hit
	target := uint32(10<<24 | 0<<16 | 0<<8 | 15)
	matched, err := engine.MatchIP(target)
	if err != nil {
		t.Fatalf("unexpected error matching IP: %v", err)
	}
	if matched.IP != target {
		t.Fatalf("expected matched IP %d, got %d", target, matched.IP)
	}
	if matched.TargetURL != "http://worker-15.internal:8080" {
		t.Fatalf("expected target URL worker-15, got %s", matched.TargetURL)
	}

	// Test miss
	missTarget := uint32(10<<24 | 0<<16 | 0<<8 | 99)
	_, err = engine.MatchIP(missTarget)
	if err != ErrRouteNotFound {
		t.Fatalf("expected ErrRouteNotFound, got %v", err)
	}
}

func TestRoutingEngine_DynamicFeatureDetection(t *testing.T) {
	profile := DetectHardwareCapabilities()

	if profile.ActiveTrack == "" {
		t.Fatal("expected non-empty active hardware track")
	}

	if profile.VectorBits != 512 && profile.VectorBits != 64 {
		t.Fatalf("expected vector bits to be 512 or 64, got %d", profile.VectorBits)
	}

	if profile.ParallelLanes != 16 && profile.ParallelLanes != 1 {
		t.Fatalf("expected parallel lanes to be 16 or 1, got %d", profile.ParallelLanes)
	}
}

func TestRoutingEngine_IPConversions(t *testing.T) {
	orig := net.ParseIP("192.168.1.100")
	asUint := IPToUint32(orig)
	back := Uint32ToIP(asUint)

	if !orig.Equal(back) {
		t.Fatalf("expected %s, got %s", orig.String(), back.String())
	}
}

func TestRoutingEngine_ConcurrentMatching(t *testing.T) {
	engine := NewRoutingEngine()

	routes := make([]RouteEntry, 64)
	for i := 0; i < 64; i++ {
		routes[i] = RouteEntry{
			IP:        uint32(172<<24 | 16<<16 | 0<<8 | (i + 1)),
			TargetURL: fmt.Sprintf("http://cluster-node-%d.prod:9000", i+1),
			Healthy:   true,
		}
	}
	engine.SetRoutes(routes)

	const workers = 32
	const iters = 1000

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				targetSlot := (workerID*37 + i) % 64
				targetIP := uint32(172<<24 | 16<<16 | 0<<8 | (targetSlot + 1))

				matched, err := engine.MatchIP(targetIP)
				if err != nil {
					t.Errorf("worker %d failed to match: %v", workerID, err)
					return
				}
				if matched.IP != targetIP {
					t.Errorf("worker %d expected %d, got %d", workerID, targetIP, matched.IP)
					return
				}
			}
		}(w)
	}

	wg.Wait()

	if total := engine.TotalEvaluations(); total != workers*iters {
		t.Fatalf("expected %d evaluations recorded, got %d", workers*iters, total)
	}
}
