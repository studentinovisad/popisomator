package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/pagination"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

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

// GetUser godoc
// @Summary Get a user by ID (admin only)
// @Tags Users
// @Produce json
// @Security CookieAuth
// @Param id path int true "User ID"
// @Success 200 {object} dto.User
// @Failure 400 {object} response.Error "invalid user ID"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "user not found"
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := service.GetUserDetails(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't get user")
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

	search, err := pagination.GetSearch(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	req := dto.ListUsersRequest{
		Limit:  limit,
		Offset: offset,
		Search: search,
	}

	query := r.URL.Query()
	if val := query.Get("role"); val != "" {
		req.Role = (*repository.UserRole)(&val)
	}

	if val := query.Get("status"); val != "" {
		req.Status = (*repository.UserStatus)(&val)
	}

	users, err := service.ListUsers(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "error fetching users")
		return
	}

	response.WriteJSON(w, http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Update a user's details, role, or status (admin only)
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
		writeServiceError(w, err, "couldn't update user")
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
