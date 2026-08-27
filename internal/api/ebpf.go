package api

import (
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/kernel"
)

type EBPFHandler struct {
	monitor *kernel.Monitor
}

func NewEBPFHandler(monitor *kernel.Monitor) *EBPFHandler {
	return &EBPFHandler{
		monitor: monitor,
	}
}

func (h *EBPFHandler) Status(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": h.monitor.Status(),
	})
}

func (h *EBPFHandler) Hooks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"hooks": h.monitor.Hooks(),
	})
}

func (h *EBPFHandler) Events(w http.ResponseWriter, r *http.Request) {
	events := make([]kernel.KernelEvent, 0)

	for {
		select {
		case event := <-h.monitor.Events():
			events = append(events, event)

		default:
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"events": events,
			})
			return
		}
	}
}
