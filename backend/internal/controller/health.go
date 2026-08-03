package controller

import (
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/db"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	if _, err := db.Queries.Healthcheck(r.Context()); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
