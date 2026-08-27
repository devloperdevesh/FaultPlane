package control

import (
	"sync"
	"time"
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
			Nodes: []TopologyNode{
				{
					ID:     "gateway",
					Type:   "gateway",
					Name:   "Gateway",
					Status: "healthy",
				},
				{
					ID:     "worker-01",
					Type:   "worker",
					Name:   "Worker-01",
					Status: "running",
				},
				{
					ID:     "checkpoint-store",
					Type:   "checkpoint",
					Name:   "Checkpoint Store",
					Status: "healthy",
				},
			},
			Connections: []TopologyConnection{
				{
					ID:     "gateway-worker-01",
					Source: "gateway",
					Target: "worker-01",
					Type:   "tcp",
					Status: "active",
				},
				{
					ID:     "worker-01-checkpoint",
					Source: "worker-01",
					Target: "checkpoint-store",
					Type:   "checkpoint-sync",
					Status: "active",
				},
			},
			UpdatedAt: time.Now().UTC(),
		},
	}
}

func (c *TopologyController) Snapshot() TopologySnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := c.snapshot

	snapshot.Nodes = append([]TopologyNode(nil), c.snapshot.Nodes...)
	snapshot.Connections = append(
		[]TopologyConnection(nil),
		c.snapshot.Connections...,
	)

	return snapshot
}

func (c *TopologyController) SetSnapshot(snapshot TopologySnapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()

	snapshot.UpdatedAt = time.Now().UTC()
	c.snapshot = snapshot
}
