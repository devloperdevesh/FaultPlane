package runtime

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/api"
	"github.com/devloperdevesh/FaultPlane/internal/models"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type WorkerRegistry struct {
	store      *api.WorkerStore
	cpuSampler *ProcessCPUSampler
}

func NewWorkerRegistry(
	store *api.WorkerStore,
	_ ...*telemetry.Registry,
) *WorkerRegistry {
	return &WorkerRegistry{
		store:      store,
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
}
