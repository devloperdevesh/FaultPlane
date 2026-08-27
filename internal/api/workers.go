package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/devloperdevesh/FaultPlane/internal/models"
)

type WorkerStore struct {
	mu      sync.RWMutex
	workers map[string]models.Worker
}

func NewWorkerStore() *WorkerStore {
	return &WorkerStore{
		workers: make(map[string]models.Worker),
	}
}

func (s *WorkerStore) Set(worker models.Worker) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.workers[worker.ID] = worker
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
