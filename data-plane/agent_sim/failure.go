// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// FailureWorkloadConfig controls the simulated connection workload.
type FailureWorkloadConfig struct {
	TargetURL       string
	RequestInterval time.Duration
	FailureAfter    int
	Timeout         time.Duration
}

// RunFailureWorkload continuously generates requests and records
// simulated connection failures and subsequent recovery.
func RunFailureWorkload(ctx context.Context, cfg FailureWorkloadConfig) error {
	if cfg.TargetURL == "" {
		return fmt.Errorf("target URL is required")
	}

	if cfg.RequestInterval <= 0 {
		cfg.RequestInterval = time.Second
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = 2 * time.Second
	}

	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	ticker := time.NewTicker(cfg.RequestInterval)
	defer ticker.Stop()

	requests := 0
	failed := false

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			requests++

			if cfg.FailureAfter > 0 && requests >= cfg.FailureAfter {
				failed = true
			}

			if failed {
				fmt.Printf(
					"agent: simulated connection failure after request %d\n",
					requests,
				)
				failed = false
				continue
			}

			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodGet,
				cfg.TargetURL,
				nil,
			)
			if err != nil {
				return fmt.Errorf("create request: %w", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("agent: request failed: %v\n", err)
				continue
			}

			resp.Body.Close()

			if resp.StatusCode >= http.StatusOK &&
				resp.StatusCode < http.StatusMultipleChoices {
				fmt.Printf("agent: workload healthy request=%d\n", requests)
			} else {
				fmt.Printf(
					"agent: workload returned status=%d request=%d\n",
					resp.StatusCode,
					requests,
				)
			}
		}
	}
}
