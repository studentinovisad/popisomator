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

func ListItems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, offset, err := pagination.GetLimitOffset(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	req := dto.ListItemsRequest{
		Limit:  limit,
		Offset: offset,
		Order:  "desc",
	}

	if val := query.Get("type_id"); val != "" {
		typeID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid type_id")
			return
		}
		req.TypeID = &typeID
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

func GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := service.GetItem(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't get item")
		return
	}

	response.WriteJSON(w, http.StatusOK, item)
}

func ConsumeItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024)

	var rawReq dto.UpdateItemRequest
	if err := json.NewDecoder(body).Decode(&rawReq); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	cleanReq := dto.UpdateItemRequest{
		ID:          id,
		Consumption: rawReq.Consumption,
	}

	item, err := service.UpdateItem(r.Context(), cleanReq)
	if err != nil {
		writeServiceError(w, err, "couldn't consume item")
		return
	}

	response.WriteJSON(w, http.StatusOK, item)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
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

	item, err := service.UpdateItem(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update item")
		return
	}

	response.WriteJSON(w, http.StatusOK, item)
}

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
