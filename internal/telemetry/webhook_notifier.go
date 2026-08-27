package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// EnterpriseDeploymentSignal represents a deployment heartbeat.
type EnterpriseDeploymentSignal struct {
	EventName         string    `json:"event_name"`
	ClientHostHash    string    `json:"client_host_hash"`
	KernelVersion     string    `json:"kernel_version"`
	Timestamp         time.Time `json:"timestamp"`
	ActiveSocketsPool uint32    `json:"active_sockets_pool"`
}

// InboundTractionSync sends deployment heartbeats to a configured telemetry endpoint.
type InboundTractionSync struct {
	TelemetryEndpoint string
	HTTPClient        *http.Client
	TransmissionCount uint64
}

// NewTractionSync creates a heartbeat sender with a bounded HTTP client.
func NewTractionSync(endpoint string) *InboundTractionSync {
	return &InboundTractionSync{
		TelemetryEndpoint: endpoint,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// EmitDeploymentHeartbeat sends one deployment heartbeat.
func (s *InboundTractionSync) EmitDeploymentHeartbeat(
	ctx context.Context,
	hostHash string,
	kernelVer string,
	socketsCount uint32,
	verificationToken string,
) (err error) {
	if s == nil || s.TelemetryEndpoint == "" {
		return fmt.Errorf("telemetry endpoint is not configured")
	}

	if s.HTTPClient == nil {
		return fmt.Errorf("telemetry HTTP client is not configured")
	}

	signal := EnterpriseDeploymentSignal{
		EventName:         "FAULTPLANE_LIVE_NODE_CONNECTED",
		ClientHostHash:    hostHash,
		KernelVersion:     kernelVer,
		Timestamp:         time.Now().UTC(),
		ActiveSocketsPool: socketsCount,
	}

	payload, err := json.Marshal(signal)
	if err != nil {
		return fmt.Errorf("marshal deployment heartbeat: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.TelemetryEndpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("create telemetry request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if verificationToken != "" {
		req.Header.Set("X-FaultPlane-Verification-Token", verificationToken)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("send deployment heartbeat: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close telemetry response body: %w", closeErr)
		}
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("telemetry endpoint returned status %s", resp.Status)
	}

	atomic.AddUint64(&s.TransmissionCount, 1)
	return nil
}
