package handlers

import (
	"net/http"

	"github.com/shandilya01/VoipalGo/internal/services"
)

type SignallingHandler struct {
	Service *services.SignallingService
}

func NewSignallingHandler() *SignallingHandler {
	return &SignallingHandler{
		Service: services.NewSignallingService(),
	}
}

func (h *SignallingHandler) HandleNewSocket(w http.ResponseWriter, r *http.Request) {
	err := h.Service.HandleNewSocketConnection(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
