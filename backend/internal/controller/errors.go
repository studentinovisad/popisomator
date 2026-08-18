package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

// writeServiceError maps a service-layer error to the appropriate HTTP status. Falls back to
// 500 with the given message for anything it doesn't recognize.
func writeServiceError(w http.ResponseWriter, err error, fallback string) {
	var pgErr *pgconn.PgError
	var valErr validator.ValidationErrors

	switch {
	case errors.Is(err, service.ErrNotFound), errors.Is(err, pgx.ErrNoRows):
		response.WriteError(w, http.StatusNotFound, "not found")
	case errors.Is(err, service.ErrInvalidReference), errors.Is(err, service.ErrNoUpdateFields), errors.Is(err, service.ErrInvalidDerivedNameFormat), errors.Is(err, service.ErrDerivedNamePropertyInUse):
		response.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.As(err, &valErr):
		response.WriteError(w, http.StatusBadRequest, "invalid request")
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		response.WriteError(w, http.StatusConflict, "already exists")
	case errors.As(err, &pgErr) && pgErr.Code == "23503":
		response.WriteError(w, http.StatusBadRequest, "invalid reference")
	default:
		response.WriteError(w, http.StatusInternalServerError, fallback)
		log.Printf("service error: %v", err)
	}

	if config.CurrentConfig.DebugMode {
		log.Printf("[DEBUG] Service error occurred. Fallback %v, Error %v", fallback, err)
	}
}
