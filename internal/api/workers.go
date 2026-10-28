package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/devloperdevesh/FaultPlane/internal/models"
)

type WorkerStore struct {
	mu      sync.RWMutex
	workers map[string]models.Worker
	events  *TelemetryStore
}

func NewWorkerStore() *WorkerStore {
	return NewWorkerStoreWithTelemetry(nil)
}

// NewWorkerStoreWithTelemetry records observed worker health transitions.
// It does not infer kernel failures or claim that a recovery was executed.
func NewWorkerStoreWithTelemetry(events *TelemetryStore) *WorkerStore {
	return &WorkerStore{workers: make(map[string]models.Worker), events: events}
}

func (s *WorkerStore) Set(worker models.Worker) {
	s.mu.Lock()
	defer s.mu.Unlock()

	previous, exists := s.workers[worker.ID]
	s.workers[worker.ID] = worker
	if s.events == nil || worker.ID == "" {
		return
	}
	status := workerHealth(worker.Status)
	previousStatus := workerHealth(previous.Status)
	eventType := ""
	if status == "failed" && (!exists || previousStatus != "failed") {
		eventType = "worker_failed"
	} else if exists && previousStatus == "failed" && status == "healthy" {
		eventType = "worker_recovered"
	}
	if eventType != "" {
		s.events.Add(TelemetryEvent{
			Type:      eventType,
			Timestamp: time.Now().UTC(),
			Metadata: map[string]string{
				"node_id": worker.ID,
				"region":  worker.Region,
				"status":  worker.Status,
				"source":  "worker_registry",
			},
		})
	}
}

func (s *WorkerStore) List() []models.Worker {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Worker, 0, len(s.workers))

	for _, worker := range s.workers {
		result = append(result, worker)
	}

	return result
}

func (s *WorkerStore) Get(id string) (models.Worker, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	worker, ok := s.workers[id]
	return worker, ok
}

func WorkersHandler(store *WorkerStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/workers")
		path = strings.Trim(path, "/")

		if path != "" {
			worker, ok := store.Get(path)
			if !ok {
				http.Error(w, "worker not found", http.StatusNotFound)
				return
			}

			writeJSON(w, http.StatusOK, worker)
			return
		}

		writeJSON(w, http.StatusOK, store.List())
	})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	payload, err := json.Marshal(value)
	if err != nil {
		http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		return
	}
}

func workerHealth(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "failed", "unhealthy", "offline", "disconnected", "dead", "error":
		return "failed"
	case "active", "healthy", "running", "ready", "recovered":
		return "healthy"
	default:
		return "unknown"
	}
}
