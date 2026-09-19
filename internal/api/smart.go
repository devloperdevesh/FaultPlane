package api

import (
	"encoding/json"
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/smart"
)

type SmartAnomalyHandler struct {
	runtime *smart.Runtime
}

func NewSmartAnomalyHandler(runtime *smart.Runtime) *SmartAnomalyHandler {
	return &SmartAnomalyHandler{
		runtime: runtime,
	}
}

func (h *SmartAnomalyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.runtime == nil {
		http.Error(w, "smart runtime unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(h.runtime.Latest()); err != nil {
		return
	}
}
