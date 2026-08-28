package runtime

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/api"
	"github.com/devloperdevesh/FaultPlane/internal/control"
	"github.com/devloperdevesh/FaultPlane/internal/models"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type WorkerRegistry struct {
	store      *api.WorkerStore
	topology   *control.TopologyController
	cpuSampler *ProcessCPUSampler
}

func NewWorkerRegistry(
	store *api.WorkerStore,
	registry *telemetry.Registry,
	topology *control.TopologyController,
) *WorkerRegistry {
	_ = registry

	return &WorkerRegistry{
		store:      store,
		topology:   topology,
		cpuSampler: NewProcessCPUSampler(),
	}
}

func (r *WorkerRegistry) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	r.publish()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.publish()
		}
	}
}

func (r *WorkerRegistry) publish() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	cpu := r.cpuSampler.Sample()

	r.store.Set(models.Worker{
		ID:     fmt.Sprintf("faultplane-runtime-%d", runtime.GOMAXPROCS(0)),
		Status: "ACTIVE",
		CPU:    cpu,
		Memory: int64(mem.Alloc),
	})

	if r.topology != nil {
		r.topology.RefreshFromWorkers(
			r.store.List(),
		)
	}
}
