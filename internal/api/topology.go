package api

import (
	"net/http"

	"github.com/devloperdevesh/FaultPlane/internal/control"
)

type TopologyHandler struct {
	controller *control.TopologyController
}

func NewTopologyHandler(controller *control.TopologyController) *TopologyHandler {
	return &TopologyHandler{
		controller: controller,
	}
}

func (h *TopologyHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		h.controller.Snapshot(),
	)
}
