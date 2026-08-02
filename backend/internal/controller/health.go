package controller

import (
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type HealthController struct {
	queries *repository.Queries
}

func NewHealthController(queries *repository.Queries) *HealthController {
	return &HealthController{queries: queries}
}

func (h *HealthController) Healthcheck(w http.ResponseWriter, r *http.Request) {
	if _, err := h.queries.Healthcheck(r.Context()); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
