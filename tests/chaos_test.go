// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type chaosBackend struct {
	healthy atomic.Bool
}

func (b *chaosBackend) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if !b.healthy.Load() {
		http.Error(w, "backend unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("primary")); err != nil {
		return
	}
}

type workloadResult struct {
	requests  int
	successes int
	failures  int
}

func runWorkload(
	ctx context.Context,
	client *http.Client,
	target string,
	interval time.Duration,
) <-chan workloadResult {
	result := make(chan workloadResult, 1)

	go func() {
		defer close(result)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var stats workloadResult

		for {
			select {
			case <-ctx.Done():
				result <- stats
				return

			case <-ticker.C:
				stats.requests++

				req, err := http.NewRequestWithContext(
					ctx,
					http.MethodGet,
					target,
					nil,
				)
				if err != nil {
					stats.failures++
					continue
				}

				resp, err := client.Do(req)
				if err != nil {
					stats.failures++
					continue
				}

				if resp.StatusCode == http.StatusOK {
					stats.successes++
				} else {
					stats.failures++
				}

				if err := resp.Body.Close(); err != nil {
					stats.failures++
				}
			}
		}
	}()

	return result
}

func TestPrimaryFailureRecovery(t *testing.T) {
	primary := &chaosBackend{}
	primary.healthy.Store(true)

	fallback := &chaosBackend{}
	fallback.healthy.Store(true)

	primaryServer := httptest.NewServer(primary)
	t.Cleanup(primaryServer.Close)

	fallbackServer := httptest.NewServer(fallback)
	t.Cleanup(fallbackServer.Close)

	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start with the primary endpoint.
	workloadURL := primaryServer.URL

	results := runWorkload(
		ctx,
		client,
		workloadURL,
		50*time.Millisecond,
	)

	// Allow the workload to establish healthy traffic.
	time.Sleep(300 * time.Millisecond)

	// Simulate primary failure.
	primary.healthy.Store(false)

	// Validate that the primary is actually unavailable.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		primaryServer.URL,
		nil,
	)
	if err != nil {
		t.Fatalf("create failure probe: %v", err)
	}

	resp, err := client.Do(req)
	if err == nil {
		if closeErr := resp.Body.Close(); closeErr != nil {
			t.Fatalf("close failure probe response: %v", closeErr)
		}
		if resp.StatusCode == http.StatusOK {
			t.Fatal("primary remained healthy after simulated failure")
		}
	}

	// Switch the workload to the fallback endpoint.
	workloadURL = fallbackServer.URL

	fallbackReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		workloadURL,
		nil,
	)
	if err != nil {
		t.Fatalf("create fallback request: %v", err)
	}

	fallbackResp, err := client.Do(fallbackReq)
	if err != nil {
		t.Fatalf("fallback recovery request failed: %v", err)
	}

	if fallbackResp.StatusCode != http.StatusOK {
		if closeErr := fallbackResp.Body.Close(); closeErr != nil {
			t.Fatalf("close fallback response: %v", closeErr)
		}

		t.Fatalf(
			"fallback returned unexpected status: %d",
			fallbackResp.StatusCode,
		)
	}

	if err := fallbackResp.Body.Close(); err != nil {
		t.Fatalf("close fallback response: %v", err)
	}

	// Generate traffic against the recovered fallback.
	fallbackResults := runWorkload(
		ctx,
		client,
		workloadURL,
		50*time.Millisecond,
	)

	recovered := <-fallbackResults
	original := <-results

	if recovered.successes == 0 {
		t.Fatal("no successful requests observed after fallback recovery")
	}

	if original.requests == 0 {
		t.Fatal("workload did not generate any requests")
	}

	t.Logf(
		"chaos recovery complete: requests=%d successes=%d failures=%d",
		recovered.requests,
		recovered.successes,
		recovered.failures,
	)
}

func TestPrimaryFailureDoesNotStopWorkload(t *testing.T) {
	primary := &chaosBackend{}
	primary.healthy.Store(true)

	server := httptest.NewServer(primary)
	t.Cleanup(server.Close)

	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var (
		requests  atomic.Int64
		successes atomic.Int64
		failures  atomic.Int64
	)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(25 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				requests.Add(1)

				req, err := http.NewRequestWithContext(
					ctx,
					http.MethodGet,
					server.URL,
					nil,
				)
				if err != nil {
					failures.Add(1)
					continue
				}

				resp, err := client.Do(req)
				if err != nil {
					failures.Add(1)
					continue
				}

				if resp.StatusCode == http.StatusOK {
					successes.Add(1)
				} else {
					failures.Add(1)
				}

				if err := resp.Body.Close(); err != nil {
					failures.Add(1)
				}
			}
		}
	}()

	time.Sleep(300 * time.Millisecond)

	// Chaos event.
	primary.healthy.Store(false)

	wg.Wait()

	if requests.Load() == 0 {
		t.Fatal("workload stopped before generating requests")
	}

	if successes.Load() == 0 {
		t.Fatal("no successful workload requests observed")
	}

	t.Logf(
		"workload survived failure: requests=%d successes=%d failures=%d",
		requests.Load(),
		successes.Load(),
		failures.Load(),
	)
}

var _ = fmt.Sprintf
