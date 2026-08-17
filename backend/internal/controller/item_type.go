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

func ListItemTypes(w http.ResponseWriter, r *http.Request) {
	limit, err := paginationValue(r, "limit", defaultPageSize, minimumPageSize, maximumPageSize)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit")
		return
	}

	offset, err := paginationValue(r, "offset", 0, 0, 0)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid offset")
		return
	}

	itemTypes, err := service.ListItemTypes(r.Context(), limit, offset)
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

func UpdateItemType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid type id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.UpdateItemTypeRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.ID = id

	itemType, err := service.UpdateItemType(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update type")
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

func AddItemTypeProperty(w http.ResponseWriter, r *http.Request) {
	typeId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid type id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.AddUpdateItemTypePropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.TypeID = typeId

	typeProp, err := service.AddItemTypeProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't add type property")
		return
	}

	response.WriteJSON(w, http.StatusOK, typeProp)
}

func UpdateItemTypeProperty(w http.ResponseWriter, r *http.Request) {
	typeId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid type id")
		return
	}

	propId, err := strconv.ParseInt(r.PathValue("prop_id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid property id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.AddUpdateItemTypePropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.TypeID = typeId
	req.PropertyID = propId

	typeProp, err := service.UpdateItemTypeProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update type property")
		return
	}

	response.WriteJSON(w, http.StatusOK, typeProp)
}

func RemoveItemTypeProperty(w http.ResponseWriter, r *http.Request) {
	typeId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid type id")
		return
	}

	propId, err := strconv.ParseInt(r.PathValue("prop_id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid property id")
		return
	}

	if err := service.RemoveItemTypeProperty(r.Context(), typeId, propId); err != nil {
		writeServiceError(w, err, "couldn't remove type property")
		return
	}

	w.WriteHeader(http.StatusOK)
}
