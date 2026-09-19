package telemetry

import "time"

// Collector connects runtime events
// with metrics registry.
type Collector struct {
	registry *Registry
}

// NewCollector creates telemetry collector.
func NewCollector(
	registry *Registry,
) *Collector {

	return &Collector{
		registry: registry,
	}
}

// RecordRequest records incoming request.
func (c *Collector) RecordRequest() {

	c.registry.IncRequests()
}

// RecordCheckpoint records checkpoint event.

// RecordRequestLatency records completed HTTP request latency.
func (c *Collector) RecordRequestLatency(duration time.Duration) {
	if duration < 0 {
		duration = 0
	}

	c.registry.RecordLatency(uint64(duration.Microseconds()) / 1000)
}

func (c *Collector) RecordCheckpoint() {

	c.registry.RecordCheckpoint()
}

// RecordRecovery records recovery event.
func (c *Collector) RecordRecovery(
	duration time.Duration,
) {

	c.registry.RecordRecovery()

	c.registry.RecordLatency(
		uint64(duration.Milliseconds()),
	)
}
