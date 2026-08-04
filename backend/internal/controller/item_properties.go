package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func GetAllProperties(w http.ResponseWriter, r *http.Request) {
	props, err := service.GetAllProperties(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get properties")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(props)
}

func GetProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	prop, err := service.GetPropertyByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't get property")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prop)
}

func CreateProperty(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.CreatePropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	prop, err := service.CreateProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create property")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prop)
}

func UpdateProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.UpdatePropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.ID = id

	prop, err := service.UpdateProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update property")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prop)
}

func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	if err := service.DeleteProperty(r.Context(), id); err != nil {
		writeServiceError(w, err, "couldn't delete property")
		return
	}

	w.WriteHeader(http.StatusOK)
}
