package control

import (
	"fmt"
	"sync"

	"github.com/devloperdevesh/FaultPlane/internal/storage"
	"github.com/devloperdevesh/FaultPlane/internal/telemetry"
)

type Controller struct {
	mu        sync.RWMutex
	workflows map[string]*Workflow
	storage   storage.Store
	telemetry *telemetry.Collector
	runtime   RuntimeExecutor
}

func NewController(
	store storage.Store,
	collector *telemetry.Collector,
	runtime RuntimeExecutor,
) *Controller {
	return &Controller{
		workflows: make(map[string]*Workflow),
		storage:   store,
		telemetry: collector,
		runtime:   runtime,
	}
}

// SetRuntime attaches the live runtime executor to the control plane.
func (c *Controller) SetRuntime(runtime RuntimeExecutor) error {
	if c == nil {
		return fmt.Errorf("set runtime: controller is nil")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if runtime == nil {
		return fmt.Errorf("set runtime: runtime is nil")
	}

	c.runtime = runtime
	return nil
}

func (c *Controller) Register(workflow *Workflow) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.workflows[workflow.ID]; exists {
		return ErrWorkflowExists
	}

	c.workflows[workflow.ID] = workflow
	return nil
}

func (c *Controller) Get(id string) (*Workflow, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	workflow, ok := c.workflows[id]
	if !ok {
		return nil, ErrWorkflowNotFound
	}

	return workflow, nil
}

func (c *Controller) Remove(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.workflows[id]; !ok {
		return ErrWorkflowNotFound
	}

	delete(c.workflows, id)
	return nil
}
