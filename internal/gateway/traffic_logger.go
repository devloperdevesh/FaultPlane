// SPDX-License-Identifier: Apache-2.0
package gateway

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrTrafficLoggerDisabled = errors.New(
		"faultplane [gateway]: traffic logging is disabled",
	)
	ErrTrafficLoggerInvalid = errors.New(
		"faultplane [gateway]: traffic logger is not initialized",
	)
)

// ActiveTractionFrame contains process-local traffic counters.
type ActiveTractionFrame struct {
	TotalProcessedBytes uint64    `json:"total_processed_bytes"`
	LivePacketCounter   uint64    `json:"live_packet_counter"`
	LastTelemetrySync   time.Time `json:"last_telemetry_sync"`
}

// ProductionTrafficLogger provides concurrent-safe audit logging for
// gateway traffic. The logger does not claim that a log file is immutable;
// filesystem permissions and storage policy remain external concerns.
type ProductionTrafficLogger struct {
	MetricsFrame      *ActiveTractionFrame
	LogStorageAddress string
	EnforcementActive uint32

	mu   sync.Mutex
	file *os.File
}

// NewTrafficLogger initializes a production traffic logger.
func NewTrafficLogger(logPath string) *ProductionTrafficLogger {
	return &ProductionTrafficLogger{
		MetricsFrame: &ActiveTractionFrame{
			LastTelemetrySync: time.Now(),
		},
		LogStorageAddress: logPath,
		EnforcementActive: 1,
	}
}

// Close releases the underlying audit log file.
func (l *ProductionTrafficLogger) Close() error {
	if l == nil {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		return nil
	}

	err := l.file.Close()
	l.file = nil
	return err
}

// LogActiveEnterpriseTraffic records a verified traffic event.
func (l *ProductionTrafficLogger) LogActiveEnterpriseTraffic(
	bytesReceived uint64,
) error {
	if l == nil || l.MetricsFrame == nil {
		return ErrTrafficLoggerInvalid
	}

	if atomic.LoadUint32(&l.EnforcementActive) == 0 {
		return ErrTrafficLoggerDisabled
	}

	if l.LogStorageAddress == "" {
		return ErrTrafficLoggerInvalid
	}

	packets := atomic.AddUint64(
		&l.MetricsFrame.LivePacketCounter,
		1,
	)
	totalBytes := atomic.AddUint64(
		&l.MetricsFrame.TotalProcessedBytes,
		bytesReceived,
	)

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		file, err := os.OpenFile(
			l.LogStorageAddress,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return fmt.Errorf(
				"faultplane [gateway]: open traffic log: %w",
				err,
			)
		}

		l.file = file
	}

	now := time.Now()
	l.MetricsFrame.LastTelemetrySync = now

	_, err := fmt.Fprintf(
		l.file,
		"[%s] FAULTPLANE TRAFFIC: packets=%d total_bytes=%d\n",
		now.Format(time.RFC3339),
		packets,
		totalBytes,
	)
	if err != nil {
		return fmt.Errorf(
			"faultplane [gateway]: write traffic log: %w",
			err,
		)
	}

	return nil
}

// Snapshot returns the current traffic counters.
func (l *ProductionTrafficLogger) Snapshot() ActiveTractionFrame {
	if l == nil || l.MetricsFrame == nil {
		return ActiveTractionFrame{}
	}

	return ActiveTractionFrame{
		TotalProcessedBytes: atomic.LoadUint64(
			&l.MetricsFrame.TotalProcessedBytes,
		),
		LivePacketCounter: atomic.LoadUint64(
			&l.MetricsFrame.LivePacketCounter,
		),
		LastTelemetrySync: l.MetricsFrame.LastTelemetrySync,
	}
}
