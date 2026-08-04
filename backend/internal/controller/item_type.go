package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func GetAllItemTypes(w http.ResponseWriter, r *http.Request) {
	itemTypes, err := service.GetAllItemTypes(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get types")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemTypes)
}

func GetItemType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid type id")
		return
	}

	itemType, err := service.GetItemType(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't get type")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemType)
}

func CreateItemType(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.CreateItemTypeRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	itemType, err := service.CreateItemType(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create item type")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemType)
}

func DeleteItemType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid type id")
		return
	}

	if err := service.DeleteItemType(r.Context(), id); err != nil {
		writeServiceError(w, err, "couldn't delete type")
		return
	}

	w.WriteHeader(http.StatusOK)
}
