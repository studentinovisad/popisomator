package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024)

	var req dto.LoginRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	token, err := service.Login(r.Context(), req)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	cookie := http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "session",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024)

	var req dto.CreateUserRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if _, err := service.GetUserByEmail(r.Context(), req.Email); err == nil {
		response.WriteError(w, http.StatusConflict, "user with this email already exists")
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		response.WriteError(w, http.StatusInternalServerError, "error checking user existence")
		return
	}

	user, err := service.CreateUser(r.Context(), req)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			response.WriteError(w, http.StatusBadRequest, "validation failed")
			return
		}

		response.WriteError(w, http.StatusInternalServerError, "error creating user")
		log.Printf("error creating user: %v", err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, user)
}
