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

// Login godoc
// @Summary Log in and receive a session cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Credentials"
// @Success 200 "Session cookie set"
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "invalid credentials"
// @Router /auth/login [post]
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

// Logout godoc
// @Summary Log out and clear the session cookie
// @Tags Auth
// @Success 200 "Session cookie cleared"
// @Router /auth/logout [post]
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

func doCreateUser(w http.ResponseWriter, r *http.Request, req dto.CreateUserRequest) {
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

// Register godoc
// @Summary Self-register a new user account (created with role "user", status "requested")
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.RegistrationRequest true "Registration details"
// @Success 201 {object} dto.User
// @Failure 400 {object} response.Error "invalid request / validation failed"
// @Failure 409 {object} response.Error "user with this email already exists"
// @Router /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024)

	var req dto.RegistrationRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	createReq := dto.CreateUserRequest{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
		Status:   "requested",
	}

	doCreateUser(w, r, createReq)
}

// CreateUser godoc
// @Summary Create a new user (admin only, created active)
// @Tags Users
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body dto.CreateUserRequest true "User to create"
// @Success 201 {object} dto.User
// @Failure 400 {object} response.Error "invalid request / validation failed"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 409 {object} response.Error "user with this email already exists"
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024)

	var req dto.CreateUserRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.Status = "active"

	doCreateUser(w, r, req)
}
