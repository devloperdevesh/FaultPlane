package runtime

import (
	"testing"

	"github.com/devloperdevesh/FaultPlane/internal/api"
	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

func TestWorkerRegistryPublishesOnlyConfiguredRegion(t *testing.T) {
	store := api.NewWorkerStore()
	topology := control.NewTopologyController()
	registry := NewWorkerRegistry(store, telemetry.NewRegistry(), topology)
	for _, region := range []string{"  operator-region  ", ""} {
		t.Setenv("FAULTPLANE_REGION", region)
		registry.publish()
		want := "operator-region"
		if region == "" {
			want = ""
		}
		if workers := store.List(); len(workers) != 1 || workers[0].Region != want {
			t.Fatalf("unexpected worker region: %+v", workers)
		}
		for _, node := range topology.Snapshot().Nodes {
			if node.Type == "worker" && node.Region != want {
				t.Fatalf("region not propagated to topology: %+v", node)
			}
			if node.Type != "worker" && node.Region != "" {
				t.Fatalf("invented infrastructure region: %+v", node)
			}
		}
	}
}
