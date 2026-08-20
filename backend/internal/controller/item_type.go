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

// GetItemTypeOptions godoc
// @Summary List all item types in a minified form, for dropdowns
// @Tags ItemTypes
// @Produce json
// @Security CookieAuth
// @Success 200 {array} dto.ItemTypeOption
// @Failure 401 {object} response.Error "not logged in"
// @Router /item-types/options [get]
func GetItemTypeOptions(w http.ResponseWriter, r *http.Request) {
	typeOptions, err := service.GetItemTypeOptions(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get item type options")
		return
	}

	response.WriteJSON(w, http.StatusOK, typeOptions)
}

// ListItemTypes godoc
// @Summary List item types
// @Tags ItemTypes
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Page size (default 20, max 50)"
// @Param offset query int false "Page offset (default 0)"
// @Param search query string false "Filter by name substring (max 100 chars)"
// @Success 200 {object} dto.ItemTypesPage
// @Failure 400 {object} response.Error "invalid limit/offset"
// @Failure 401 {object} response.Error "not logged in"
// @Router /item-types [get]
func ListItemTypes(w http.ResponseWriter, r *http.Request) {
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

	itemTypes, err := service.ListItemTypes(r.Context(), limit, offset, search)
	if err != nil {
		writeServiceError(w, err, "couldn't list types")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemTypes)
}

// GetItemType godoc
// @Summary Get an item type by ID
// @Tags ItemTypes
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item Type ID"
// @Success 200 {object} dto.ItemType
// @Failure 400 {object} response.Error "invalid type id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 404 {object} response.Error "not found"
// @Router /item-types/{id} [get]
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

// CreateItemType godoc
// @Summary Create an item type (manager/admin only)
// @Tags ItemTypes
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body dto.CreateItemTypeRequest true "Item type to create"
// @Success 200 {object} dto.ItemType
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /item-types [post]
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

// UpdateItemType godoc
// @Summary Update an item type (manager/admin only)
// @Tags ItemTypes
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item Type ID"
// @Param body body dto.UpdateItemTypeRequest true "Fields to update"
// @Success 200 {object} dto.ItemType
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /item-types/{id} [patch]
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

// DeleteItemType godoc
// @Summary Delete an item type (manager/admin only)
// @Tags ItemTypes
// @Security CookieAuth
// @Param id path int true "Item Type ID"
// @Success 200
// @Failure 400 {object} response.Error "invalid type id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /item-types/{id} [delete]
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

// AddItemTypeProperty godoc
// @Summary Add a property to an item type (manager/admin only)
// @Tags ItemTypes
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item Type ID"
// @Param body body dto.AddUpdateItemTypePropertyRequest true "Property to add"
// @Success 200 {object} dto.ItemTypeProperty
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /item-types/{id}/properties [post]
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

// UpdateItemTypeProperty godoc
// @Summary Update a property on an item type (manager/admin only)
// @Tags ItemTypes
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item Type ID"
// @Param prop_id path int true "Property ID"
// @Param body body dto.AddUpdateItemTypePropertyRequest true "Property fields"
// @Success 200 {object} dto.ItemTypeProperty
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /item-types/{id}/properties/{prop_id} [patch]
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

// RemoveItemTypeProperty godoc
// @Summary Remove a property from an item type (manager/admin only)
// @Tags ItemTypes
// @Security CookieAuth
// @Param id path int true "Item Type ID"
// @Param prop_id path int true "Property ID"
// @Success 200
// @Failure 400 {object} response.Error "invalid type/property id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /item-types/{id}/properties/{prop_id} [delete]
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
