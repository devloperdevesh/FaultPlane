package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type TelemetryEvent struct {
	Type      string            `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	Value     float64           `json:"value,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type TelemetryResponse struct {
	Events []TelemetryEvent `json:"events"`
}

type LogsResponse struct {
	Logs []TelemetryEvent `json:"logs"`
}

type NetworkEventsResponse struct {
	Events []TelemetryEvent `json:"events"`
}

type TelemetryStore struct {
	mu     sync.RWMutex
	events []TelemetryEvent
}

func NewTelemetryStore() *TelemetryStore {
	return &TelemetryStore{
		events: make([]TelemetryEvent, 0),
	}
}

func (s *TelemetryStore) Add(event TelemetryEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)

	const maxEvents = 1000

	if len(s.events) > maxEvents {
		s.events = s.events[len(s.events)-maxEvents:]
	}
}

func (s *TelemetryStore) Snapshot() []TelemetryEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]TelemetryEvent, len(s.events))
	copy(result, s.events)

	return result
}

func TelemetryHandler(store *TelemetryStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		writeTelemetryJSON(w, http.StatusOK, TelemetryResponse{
			Events: store.Snapshot(),
		})
	})
}

func LogsHandler(store *TelemetryStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		writeTelemetryJSON(w, http.StatusOK, LogsResponse{
			Logs: store.Snapshot(),
		})
	})
}

func NetworkEventsHandler(store *TelemetryStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		writeTelemetryJSON(w, http.StatusOK, NetworkEventsResponse{
			Events: store.Snapshot(),
		})
	})
}

func writeTelemetryJSON(
	w http.ResponseWriter,
	status int,
	value interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
