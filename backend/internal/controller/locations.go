package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

// CreateLocation godoc
// @Summary Create a location (manager/admin only)
// @Tags Locations
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body dto.CreateLocationRequest true "Location to create"
// @Success 200 {object} dto.Location
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /locations [post]
func CreateLocation(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.CreateLocationRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	location, err := service.CreateLocation(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create location")
		return
	}

	response.WriteJSON(w, http.StatusOK, location)
}

// GetLocationOptionsFlat godoc
// @Summary List all location options flat
// @Tags Locations
// @Produce json
// @Security CookieAuth
// @Success 200 {array} dto.LocationOption
// @Failure 401 {object} response.Error "not logged in"
// @Router /locations/flat [get]
func GetLocationOptionsFlat(w http.ResponseWriter, r *http.Request) {
	locations, err := service.GetLocationOptionsFlat(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get location options")
		return
	}

	response.WriteJSON(w, http.StatusOK, locations)
}

// GetLocationOptions godoc
// @Summary List all locations in a minified form
// @Tags Locations
// @Produce json
// @Security CookieAuth
// @Success 200 {array} dto.LocationOption
// @Failure 401 {object} response.Error "not logged in"
// @Router /locations [get]
func GetLocationOptions(w http.ResponseWriter, r *http.Request) {
	locations, err := service.GetLocationOptions(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get location options")
		return
	}

	response.WriteJSON(w, http.StatusOK, locations)
}

// GetLocation godoc
// @Summary Get a location by ID
// @Tags Locations
// @Produce json
// @Security CookieAuth
// @Param id path int true "Location ID"
// @Success 200 {object} dto.Location
// @Failure 400 {object} response.Error "invalid location id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 404 {object} response.Error "not found"
// @Router /locations/{id} [get]
func GetLocation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid property id")
		return
	}

	location, err := service.GetLocationByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't get location")
		return
	}

	response.WriteJSON(w, http.StatusOK, location)
}

// UpdateLocation godoc
// @Summary Update a location (manager/admin only)
// @Tags Locations
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Location ID"
// @Param body body dto.UpdateLocationRequest true "Fields to update"
// @Success 200 {object} dto.Location
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /locations/{id} [patch]
func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid location id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.UpdateLocationRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.ID = id

	location, err := service.UpdateLocation(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update location")
		return
	}

	response.WriteJSON(w, http.StatusOK, location)
}

// DeleteLocation godoc
// @Summary Delete a location (manager/admin only)
// @Tags Locations
// @Security CookieAuth
// @Param id path int true "Location ID"
// @Success 200
// @Failure 400 {object} response.Error "invalid location id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /locations/{id} [delete]
func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid location id")
		return
	}

	if err := service.DeleteLocation(r.Context(), id); err != nil {
		writeServiceError(w, err, "couldn't delete location")
		return
	}

	w.WriteHeader(http.StatusOK)
}
