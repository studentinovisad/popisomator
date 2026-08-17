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

func GetAllProperties(w http.ResponseWriter, r *http.Request) {
	props, err := service.GetAllProperties(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get properties")
		return
	}

	response.WriteJSON(w, http.StatusOK, props)
}

func ListProperties(w http.ResponseWriter, r *http.Request) {
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

	properties, err := service.ListProperties(r.Context(), limit, offset)
	if err != nil {
		writeServiceError(w, err, "couldn't get properties")
		return
	}

	response.WriteJSON(w, http.StatusOK, properties)
}

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
