package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

// writeServiceError maps a service-layer error to the appropriate HTTP status. Falls back to
// 500 with the given message for anything it doesn't recognize.
func writeServiceError(w http.ResponseWriter, err error, fallback string) {
	var pgErr *pgconn.PgError
	var valErr validator.ValidationErrors

	switch {
	case errors.Is(err, service.ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, service.ErrInvalidReference):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.As(err, &valErr):
		http.Error(w, "invalid request", http.StatusBadRequest)
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		http.Error(w, "already exists", http.StatusConflict)
	case errors.As(err, &pgErr) && pgErr.Code == "23503":
		http.Error(w, "invalid reference", http.StatusBadRequest)
	default:
		http.Error(w, fallback, http.StatusInternalServerError)
		log.Printf("service error: %v", err)
	}
}
