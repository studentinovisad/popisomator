package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/pagination"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

// ListItems godoc
// @Summary List items
// @Tags Items
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Page size (default 20, max 50)"
// @Param offset query int false "Page offset (default 0)"
// @Param search query string false "Filter by derived item name substring (max 100 chars)"
// @Param type_id query int false "Filter by item type ID"
// @Param property.{id} query string false "Exact value filter for an item-type property (max 100 chars)"
// @Param consumption query []string false "Filter by consumption status (comma-separated)" collectionFormat(csv) Enums(not_consumed, partially_consumed, fully_consumed, damaged)
// @Param created_from query string false "Filter by creation time, RFC3339"
// @Param created_to query string false "Filter by creation time, RFC3339"
// @Param order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} dto.ItemsPage
// @Failure 400 {object} response.Error "invalid query parameters"
// @Failure 401 {object} response.Error "not logged in"
// @Router /items [get]
func ListItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}
	query := r.URL.Query()

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

	req := dto.ListItemsRequest{
		Limit:    limit,
		Offset:   offset,
		Order:    "desc",
		Search:   search,
		ViewerID: userID,
	}

	if val := query.Get("type_id"); val != "" {
		typeID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid type_id")
			return
		}
		req.TypeID = &typeID
	}

	for key, values := range query {
		propertyIDText, ok := strings.CutPrefix(key, "property.")
		if !ok {
			continue
		}

		propertyID, err := strconv.ParseInt(propertyIDText, 10, 64)
		if err != nil || propertyID < 1 || len(values) != 1 || len(values[0]) > 100 {
			response.WriteError(w, http.StatusBadRequest, "invalid property filter")
			return
		}
		if req.PropertyFilters == nil {
			req.PropertyFilters = make(map[int64]json.RawMessage)
		}
		req.PropertyFilters[propertyID] = json.RawMessage(values[0])
	}

	if vals, ok := query["consumption"]; ok {
		for _, v := range vals {
			for _, part := range strings.Split(v, ",") {
				if part == "" {
					continue
				}
				req.Consumption = append(req.Consumption, repository.ConsumptionStatus(part))
			}
		}
	}

	if val := query.Get("created_from"); val != "" {
		parsed_time, err := time.Parse(time.RFC3339, val)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid created_from")
			return
		}
		req.CreatedFrom = &parsed_time
	}

	if val := query.Get("created_to"); val != "" {
		parsed_time, err := time.Parse(time.RFC3339, val)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid created_to")
			return
		}
		req.CreatedTo = &parsed_time
	}

	if val := query.Get("order"); val != "" {
		req.Order = val
	}

	result, err := service.ListItems(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't list items")
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
}

// GetItem godoc
// @Summary Get an item by ID
// @Tags Items
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item ID"
// @Success 200 {object} dto.Item
// @Failure 400 {object} response.Error "invalid item id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 404 {object} response.Error "not found"
// @Router /items/{id} [get]
func GetItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := service.GetItem(r.Context(), id, userID)
	if err != nil {
		writeServiceError(w, err, "couldn't get item")
		return
	}

	response.WriteJSON(w, http.StatusOK, item)
}

// ConsumeItem godoc
// @Summary Update an item's consumption status
// @Tags Items
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item ID"
// @Param body body dto.UpdateItemRequest true "Consumption status"
// @Success 200 {object} dto.Item
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "not approved to consume item"
// @Failure 404 {object} response.Error "not found"
// @Router /items/{id}/consume [post]
func ConsumeItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	// Check for approval before consuming item
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}

	userDetails, err := service.GetUserDetails(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "error fetching role")
		return
	}
	if userDetails.Role == "user" {
		isApproved, err := service.CheckItemApproval(r.Context(), userID, itemID)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, "error checking for approval")
			return
		} else if !isApproved {
			response.WriteError(w, http.StatusForbidden, "not approved to consume item")
			return
		}
	}

	body := http.MaxBytesReader(w, r.Body, 1024)

	var rawReq dto.UpdateItemRequest
	if err := json.NewDecoder(body).Decode(&rawReq); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	cleanReq := dto.UpdateItemRequest{
		ID:          itemID,
		Consumption: rawReq.Consumption,
		ViewerID:    userID,
	}

	item, err := service.UpdateItem(r.Context(), cleanReq)
	if err != nil {
		writeServiceError(w, err, "couldn't consume item")
		return
	}

	response.WriteJSON(w, http.StatusOK, item)
}

// UpdateItem godoc
// @Summary Update an item's type (manager/admin only)
// @Tags Items
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item ID"
// @Param body body dto.UpdateItemRequest true "Fields to update"
// @Success 200 {object} dto.Item
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /items/{id} [patch]
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024)

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	var req dto.UpdateItemRequest
	if err := decoder.Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.ID = id
	req.Consumption = nil
	req.ViewerID = userID

	item, err := service.UpdateItem(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update item")
		return
	}

	response.WriteJSON(w, http.StatusOK, item)
}

// CreateItem godoc
// @Summary Create one or more items of a type (manager/admin only)
// @Tags Items
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body dto.CreateItemRequest true "Item(s) to create"
// @Success 200 {array} dto.Item
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /items [post]
func CreateItem(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*64)

	var req dto.CreateItemRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	items, err := service.CreateItem(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create items")
		return
	}

	response.WriteJSON(w, http.StatusOK, items)
}

// DeleteItem godoc
// @Summary Delete an item (manager/admin only)
// @Tags Items
// @Security CookieAuth
// @Param id path int true "Item ID"
// @Success 200
// @Failure 400 {object} response.Error "invalid item id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /items/{id} [delete]
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := service.DeleteItem(r.Context(), id); err != nil {
		writeServiceError(w, err, "couldn't delete item")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AddItemProperty godoc
// @Summary Add a property value to an item (manager/admin only)
// @Tags Items
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item ID"
// @Param body body dto.AddUpdateItemPropertyRequest true "Property to add"
// @Success 200 {object} dto.ItemProperty
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /items/{id}/properties [post]
func AddItemProperty(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.AddUpdateItemPropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.ItemID = itemId

	itemProp, err := service.AddItemProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't add item property")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemProp)
}

// UpdateItemProperty godoc
// @Summary Update a property value on an item (manager/admin only)
// @Tags Items
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param id path int true "Item ID"
// @Param prop_id path int true "Property ID"
// @Param body body dto.AddUpdateItemPropertyRequest true "Property value"
// @Success 200 {object} dto.ItemProperty
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /items/{id}/properties/{prop_id} [put]
func UpdateItemProperty(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	propId, err := strconv.ParseInt(r.PathValue("prop_id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid property id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.AddUpdateItemPropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.ItemID = itemId
	req.PropertyID = propId

	itemProp, err := service.UpdateItemProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update property of item")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemProp)
}

// RemoveItemProperty godoc
// @Summary Remove a property value from an item (manager/admin only)
// @Tags Items
// @Security CookieAuth
// @Param id path int true "Item ID"
// @Param prop_id path int true "Property ID"
// @Success 200
// @Failure 400 {object} response.Error "invalid item/property id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /items/{id}/properties/{prop_id} [delete]
func RemoveItemProperty(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	propId, err := strconv.ParseInt(r.PathValue("prop_id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid property id")
		return
	}

	if err := service.RemoveItemProperty(r.Context(), itemId, propId); err != nil {
		writeServiceError(w, err, "couldn't remove item property")
		return
	}

	w.WriteHeader(http.StatusOK)
}
