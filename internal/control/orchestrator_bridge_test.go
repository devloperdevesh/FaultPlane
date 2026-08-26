package control

import "testing"

func TestOrchestratorBridgeHealthyProvider(t *testing.T) {
	bridge := NewOrchestratorBridge(RoutingPolicy{
		MaxLatencyMs:       100,
		MaxErrorRate:       0.05,
		MinTokensPerSecond: 10,
	})

	decision := bridge.Evaluate(ProviderSignal{
		Provider:        ProviderVLLM,
		Healthy:         true,
		LatencyMs:       35,
		ErrorRate:       0.01,
		TokensPerSecond: 50,
	})

	if decision.State != RoutingHealthy {
		t.Fatalf("expected healthy state, got %s", decision.State)
	}
}

func TestOrchestratorBridgeFailoverOnUnhealthyProvider(t *testing.T) {
	bridge := NewOrchestratorBridge(RoutingPolicy{})

	decision := bridge.Evaluate(ProviderSignal{
		Provider: ProviderOpenAI,
		Healthy:  false,
	})

	if decision.State != RoutingFailover {
		t.Fatalf("expected failover state, got %s", decision.State)
	}
}

func TestOrchestratorBridgeLatencyDegradation(t *testing.T) {
	bridge := NewOrchestratorBridge(RoutingPolicy{
		MaxLatencyMs: 100,
	})

	decision := bridge.Evaluate(ProviderSignal{
		Provider:  ProviderAnthropic,
		Healthy:   true,
		LatencyMs: 250,
	})

	if decision.State != RoutingDegraded {
		t.Fatalf("expected degraded state, got %s", decision.State)
	}
}

func TestRoutingPolicyValidation(t *testing.T) {
	policy := RoutingPolicy{
		MaxLatencyMs: -1,
	}

	if err := policy.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}
