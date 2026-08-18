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

func GetPropertyOptions(w http.ResponseWriter, r *http.Request) {
	props, err := service.GetPropertyOptions(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get property options")
		return
	}

	response.WriteJSON(w, http.StatusOK, props)
}

func ListProperties(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := pagination.GetLimitOffset(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	properties, err := service.ListProperties(r.Context(), limit, offset)
	if err != nil {
		writeServiceError(w, err, "couldn't list properties")
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
