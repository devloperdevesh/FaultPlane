package control

import (
	"sync"

	"github.com/devloperdevesh/agentmesh/internal/storage"
	"github.com/devloperdevesh/agentmesh/internal/telemetry"
)

// Controller manages agent workflow lifecycle,
// state tracking and recovery operations.
type Controller struct {

	// protects concurrent workflow access
	mu sync.RWMutex

	// workflow registry
	workflows map[string]*Workflow

	// persistent state storage
	storage storage.Store

	// telemetry collector
	telemetry *telemetry.Collector
}

// NewController creates a new control plane controller.
func NewController(
	store storage.Store,
	collector *telemetry.Collector,
) *Controller {

	return &Controller{

		workflows: make(map[string]*Workflow),

		storage: store,

		telemetry: collector,
	}
}

// Register adds a new workflow into controller state.
func (c *Controller) Register(
	id string,
	workflow *Workflow,
) {

	c.mu.Lock()
	defer c.mu.Unlock()

	c.workflows[id] = workflow
}

// Get retrieves workflow by ID.
func (c *Controller) Get(
	id string,
) (*Workflow, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()

	workflow, exists := c.workflows[id]

	if !exists {
		return nil, ErrWorkflowNotFound
	}

	return workflow, nil
}

// Remove deletes workflow from controller.
func (c *Controller) Remove(
	id string,
) {

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.workflows, id)
}
