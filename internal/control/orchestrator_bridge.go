package control

import "fmt"

// Provider identifies an inference/runtime provider.
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderVLLM      Provider = "vllm"
	ProviderTriton    Provider = "triton"
)

// RoutingState represents the current routing decision state.
type RoutingState string

const (
	RoutingHealthy  RoutingState = "healthy"
	RoutingDegraded RoutingState = "degraded"
	RoutingFailover RoutingState = "failover"
	RoutingDisabled RoutingState = "disabled"
)

// ProviderSignal contains runtime signals used by the router.
type ProviderSignal struct {
	Provider        Provider
	Healthy         bool
	LatencyMs       float64
	ErrorRate       float64
	TokensPerSecond float64
}

// RoutingPolicy controls when a provider should be considered unhealthy.
type RoutingPolicy struct {
	MaxLatencyMs       float64
	MaxErrorRate       float64
	MinTokensPerSecond float64
}

// RoutingDecision is the result of evaluating one provider.
type RoutingDecision struct {
	Provider Provider
	State    RoutingState
	Reason   string
}

// OrchestratorBridge evaluates runtime provider signals.
type OrchestratorBridge struct {
	policy RoutingPolicy
}

// NewOrchestratorBridge creates a routing bridge with the supplied policy.
func NewOrchestratorBridge(policy RoutingPolicy) *OrchestratorBridge {
	return &OrchestratorBridge{
		policy: policy,
	}
}

// Evaluate converts provider telemetry into a routing decision.
func (b *OrchestratorBridge) Evaluate(
	signal ProviderSignal,
) RoutingDecision {
	if signal.Provider == "" {
		return RoutingDecision{
			State:  RoutingDisabled,
			Reason: "provider is not specified",
		}
	}

	if !signal.Healthy {
		return RoutingDecision{
			Provider: signal.Provider,
			State:    RoutingFailover,
			Reason:   "provider reported unhealthy",
		}
	}

	if b.policy.MaxLatencyMs > 0 &&
		signal.LatencyMs > b.policy.MaxLatencyMs {
		return RoutingDecision{
			Provider: signal.Provider,
			State:    RoutingDegraded,
			Reason:   "latency threshold exceeded",
		}
	}

	if b.policy.MaxErrorRate > 0 &&
		signal.ErrorRate > b.policy.MaxErrorRate {
		return RoutingDecision{
			Provider: signal.Provider,
			State:    RoutingDegraded,
			Reason:   "error-rate threshold exceeded",
		}
	}

	if b.policy.MinTokensPerSecond > 0 &&
		signal.TokensPerSecond < b.policy.MinTokensPerSecond {
		return RoutingDecision{
			Provider: signal.Provider,
			State:    RoutingDegraded,
			Reason:   "token throughput below threshold",
		}
	}

	return RoutingDecision{
		Provider: signal.Provider,
		State:    RoutingHealthy,
		Reason:   "provider satisfies routing policy",
	}
}

// Validate checks whether a routing policy is usable.
func (p RoutingPolicy) Validate() error {
	if p.MaxLatencyMs < 0 {
		return fmt.Errorf("max latency cannot be negative")
	}

	if p.MaxErrorRate < 0 {
		return fmt.Errorf("max error rate cannot be negative")
	}

	if p.MinTokensPerSecond < 0 {
		return fmt.Errorf("minimum token throughput cannot be negative")
	}

	return nil
}
