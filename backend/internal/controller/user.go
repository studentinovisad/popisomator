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

func ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, err := pagination.QueryValue(r, "limit", pagination.DefaultPageSize, pagination.MinimumPageSize, pagination.MaximumPageSize)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit")
		return
	}

	offset, err := pagination.QueryValue(r, "offset", 0, 0, 0)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid offset")
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

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024)

	var req dto.UpdateRoleRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := service.UpdateUserRole(r.Context(), id, req)
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

		response.WriteError(w, http.StatusInternalServerError, "error updating role")
		log.Printf("error updating role: %v", err)
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

func DeclineRegistration(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := service.DeclineRegistration(r.Context(), id); err != nil {
		writeServiceError(w, err, "couldn't decline user registration")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ApproveRegistration(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := service.ApproveRegistration(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't approve user registration")
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}
