package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/pagination"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

// CreateProperty godoc
// @Summary Create a property definition (manager/admin only)
// @Tags Properties
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body dto.CreatePropertyRequest true "Property to create"
// @Success 200 {object} dto.Property
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /properties [post]
func CreateProperty(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.CreatePropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	prop, err := service.CreateProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create property")
		return
	}

	response.WriteJSON(w, http.StatusOK, prop)
}

// GetPropertyOptions godoc
// @Summary List all properties in a minified form, for dropdowns
// @Tags Properties
// @Produce json
// @Security CookieAuth
// @Success 200 {array} dto.PropertyOption
// @Failure 401 {object} response.Error "not logged in"
// @Router /properties/options [get]
func GetPropertyOptions(w http.ResponseWriter, r *http.Request) {
	props, err := service.GetPropertyOptions(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get property options")
		return
	}

	response.WriteJSON(w, http.StatusOK, props)
}

// ListProperties godoc
// @Summary List property definitions
// @Tags Properties
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Page size (default 20, max 50)"
// @Param offset query int false "Page offset (default 0)"
// @Param search query string false "Filter by name substring (max 100 chars)"
// @Success 200 {object} dto.PropertiesPage
// @Failure 400 {object} response.Error "invalid limit/offset"
// @Failure 401 {object} response.Error "not logged in"
// @Router /properties [get]
func ListProperties(w http.ResponseWriter, r *http.Request) {
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

	properties, err := service.ListProperties(r.Context(), limit, offset, search)
	if err != nil {
		writeServiceError(w, err, "couldn't list properties")
		return
	}

	response.WriteJSON(w, http.StatusOK, properties)
}

// GetProperty godoc
// @Summary Get a property definition by ID
// @Tags Properties
// @Produce json
// @Security CookieAuth
// @Param id path int true "Property ID"
// @Success 200 {object} dto.Property
// @Failure 400 {object} response.Error "invalid property id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 404 {object} response.Error "not found"
// @Router /properties/{id} [get]
func GetProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid property id")
		return
	}

	prop, err := service.GetPropertyByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't get property")
		return
	}

	response.WriteJSON(w, http.StatusOK, prop)
}

// UpdateProperty godoc
// @Summary Update a property definition (manager/admin only)
// @Tags Properties
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Property ID"
// @Param body body dto.UpdatePropertyRequest true "Fields to update"
// @Success 200 {object} dto.Property
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /properties/{id} [patch]
func UpdateProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid property id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.UpdatePropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.ID = id

	prop, err := service.UpdateProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update property")
		return
	}

	response.WriteJSON(w, http.StatusOK, prop)
}

// DeleteProperty godoc
// @Summary Delete a property definition (manager/admin only)
// @Tags Properties
// @Security CookieAuth
// @Param id path int true "Property ID"
// @Success 200
// @Failure 400 {object} response.Error "invalid property id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /properties/{id} [delete]
func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid property id")
		return
	}

	if err := service.DeleteProperty(r.Context(), id); err != nil {
		writeServiceError(w, err, "couldn't delete property")
		return
	}

	w.WriteHeader(http.StatusOK)
}
