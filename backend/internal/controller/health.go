package controller

import (
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/db"
)

// Healthcheck godoc
// @Summary Readiness probe, verifies database connectivity
// @Tags Operational
// @Produce plain
// @Success 200 {string} string "OK"
// @Failure 503 {string} string "database unavailable"
// @Router /health [get]
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	if _, err := db.Queries.Healthcheck(r.Context()); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
