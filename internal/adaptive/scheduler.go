package adaptive

import (
	"sync/atomic"
	"time"
)

type InvariantMetricsFrame struct {
	TailLatencyNS     int64
	ActiveTokensCount uint64
	ProviderDropFlag  uint32
}

type AdaptiveScheduler struct {
	MetricsFrame       *InvariantMetricsFrame
	TargetLatencyBound time.Duration
}

func NewAdaptiveScheduler(bound time.Duration) *AdaptiveScheduler {
	if bound <= 0 {
		bound = 2 * time.Millisecond
	}

	return &AdaptiveScheduler{
		MetricsFrame: &InvariantMetricsFrame{
			TailLatencyNS: bound.Nanoseconds(),
		},
		TargetLatencyBound: bound,
	}
}

func (s *AdaptiveScheduler) EvaluateProviderSoundness(
	currentLatency time.Duration,
) bool {
	atomic.StoreInt64(
		&s.MetricsFrame.TailLatencyNS,
		currentLatency.Nanoseconds(),
	)

	if currentLatency > s.TargetLatencyBound {
		atomic.StoreUint32(
			&s.MetricsFrame.ProviderDropFlag,
			1,
		)
		return false
	}

	atomic.StoreUint32(
		&s.MetricsFrame.ProviderDropFlag,
		0,
	)

	return true
}

func (s *AdaptiveScheduler) SetActiveTokens(count uint64) {
	atomic.StoreUint64(
		&s.MetricsFrame.ActiveTokensCount,
		count,
	)
}

func (s *AdaptiveScheduler) Snapshot() InvariantMetricsFrame {
	return InvariantMetricsFrame{
		TailLatencyNS: atomic.LoadInt64(
			&s.MetricsFrame.TailLatencyNS,
		),
		ActiveTokensCount: atomic.LoadUint64(
			&s.MetricsFrame.ActiveTokensCount,
		),
		ProviderDropFlag: atomic.LoadUint32(
			&s.MetricsFrame.ProviderDropFlag,
		),
	}
}
