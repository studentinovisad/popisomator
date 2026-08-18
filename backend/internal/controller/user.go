package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/pagination"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

const maxUserSearchLength = 100

// UserDetailsPersonal godoc
// @Summary Get the currently authenticated user's own details
// @Tags Users
// @Produce json
// @Security CookieAuth
// @Success 200 {object} dto.User
// @Failure 401 {object} response.Error "not logged in"
// @Router /users/me [get]
func UserDetailsPersonal(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}

	user, err := service.GetUserDetails(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "error fetching details")
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// ListUsers godoc
// @Summary List users (admin only)
// @Tags Users
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Page size (default 20, max 50)"
// @Param offset query int false "Page offset (default 0)"
// @Param search query string false "Filter by name/email substring (max 100 chars)"
// @Param role query string false "Filter by role" Enums(admin, manager, user)
// @Param status query string false "Filter by status" Enums(requested, active)
// @Success 200 {object} dto.UsersPage
// @Failure 400 {object} response.Error "invalid limit/offset/search/role/status"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /users [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := pagination.GetLimitOffset(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	if len(search) > maxUserSearchLength {
		response.WriteError(w, http.StatusBadRequest, "search is too long")
		return
	}

	role := strings.TrimSpace(r.URL.Query().Get("role"))
	if role != "" && role != "admin" && role != "manager" && role != "user" {
		response.WriteError(w, http.StatusBadRequest, "invalid role")
		return
	}

	status := strings.TrimSpace(r.URL.Query().Get("status"))
	if status != "" && status != "requested" && status != "active" {
		response.WriteError(w, http.StatusBadRequest, "invalid status")
		return
	}

	users, err := service.ListUsers(r.Context(), dto.ListUsersRequest{
		Limit:  limit,
		Offset: offset,
		Search: search,
		Role:   role,
		Status: status,
	})
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "error fetching users")
		log.Printf("error fetching users: %v", err)
		return
	}

	response.WriteJSON(w, http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Update a user's role/status (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "User ID"
// @Param body body dto.UpdateUserRequest true "Fields to update"
// @Success 200 {object} dto.User
// @Failure 400 {object} response.Error "invalid request / validation failed"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "user not found"
// @Router /users/{id} [patch]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024)

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := service.UpdateUser(r.Context(), id, req)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			response.WriteError(w, http.StatusBadRequest, "validation failed")
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "user not found")
			return
		}

		response.WriteError(w, http.StatusInternalServerError, "error updating user")
		log.Printf("error updating user: %v", err)
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user (admin only)
// @Tags Users
// @Security CookieAuth
// @Param id path int true "User ID"
// @Success 200
// @Failure 400 {object} response.Error "invalid user id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := service.DeleteUser(r.Context(), id); err != nil {
		writeServiceError(w, err, "couldn't delete user")
		return
	}

	w.WriteHeader(http.StatusOK)
}
