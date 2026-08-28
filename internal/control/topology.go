package control

import (
	"sync"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/models"
)

type TopologyNode struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type TopologyConnection struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type TopologySnapshot struct {
	Nodes       []TopologyNode       `json:"nodes"`
	Connections []TopologyConnection `json:"connections"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type TopologyController struct {
	mu       sync.RWMutex
	snapshot TopologySnapshot
}

func NewTopologyController() *TopologyController {
	return &TopologyController{
		snapshot: TopologySnapshot{
			Nodes:       []TopologyNode{},
			Connections: []TopologyConnection{},
			UpdatedAt:   time.Now().UTC(),
		},
	}
}

func (c *TopologyController) Snapshot() TopologySnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := c.snapshot

	snapshot.Nodes = append(
		[]TopologyNode(nil),
		c.snapshot.Nodes...,
	)

	snapshot.Connections = append(
		[]TopologyConnection(nil),
		c.snapshot.Connections...,
	)

	return snapshot
}

func (c *TopologyController) SetSnapshot(
	snapshot TopologySnapshot,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	snapshot.UpdatedAt = time.Now().UTC()
	c.snapshot = snapshot
}

func (c *TopologyController) RefreshFromWorkers(
	workers []models.Worker,
) {
	nodes := make(
		[]TopologyNode,
		0,
		len(workers)+2,
	)

	connections := make(
		[]TopologyConnection,
		0,
		len(workers)*2,
	)

	nodes = append(
		nodes,
		TopologyNode{
			ID:     "gateway",
			Type:   "gateway",
			Name:   "Gateway",
			Status: "healthy",
		},
	)

	for _, worker := range workers {
		nodes = append(
			nodes,
			TopologyNode{
				ID:     worker.ID,
				Type:   "worker",
				Name:   worker.ID,
				Status: worker.Status,
			},
		)

		connections = append(
			connections,
			TopologyConnection{
				ID:     "gateway-" + worker.ID,
				Source: "gateway",
				Target: worker.ID,
				Type:   "runtime",
				Status: "active",
			},
		)
	}

	nodes = append(
		nodes,
		TopologyNode{
			ID:     "checkpoint-store",
			Type:   "checkpoint",
			Name:   "Checkpoint Store",
			Status: "healthy",
		},
	)

	for _, worker := range workers {
		connections = append(
			connections,
			TopologyConnection{
				ID:     worker.ID + "-checkpoint",
				Source: worker.ID,
				Target: "checkpoint-store",
				Type:   "checkpoint-sync",
				Status: "active",
			},
		)
	}

	c.SetSnapshot(
		TopologySnapshot{
			Nodes:       nodes,
			Connections: connections,
		},
	)
}
