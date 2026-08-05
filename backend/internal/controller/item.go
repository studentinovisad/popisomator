package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

const defaultItemPageSize int32 = 20
const maxItemPageSize int32 = 100

// paginationValue parses a query param as an int32, falling back to fallback when absent and
// erroring when unparseable or outside [minimum, maximum] (maximum <= 0 disables the upper bound).
func paginationValue(r *http.Request, key string, fallback, minimum, maximum int32) (int32, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}

	parsedValue := int32(parsed)
	if parsedValue < minimum || (maximum > 0 && parsedValue > maximum) {
		return 0, strconv.ErrSyntax
	}

	return parsedValue, nil
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, err := paginationValue(r, "limit", defaultItemPageSize, 1, maxItemPageSize)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit")
		return
	}

	offset, err := paginationValue(r, "offset", 0, 0, 0)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid offset")
		return
	}

	req := dto.ListItemsRequest{
		Limit:  limit,
		Offset: offset,
		Order:  "desc",
	}

	if v := q.Get("type_id"); v != "" {
		typeID, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid type_id")
			return
		}
		req.TypeID = &typeID
	}

	if vals, ok := q["consumption"]; ok {
		for _, v := range vals {
			for _, part := range strings.Split(v, ",") {
				if part == "" {
					continue
				}
				req.Consumption = append(req.Consumption, repository.ConsumptionStatus(part))
			}
		}
	}

	if v := q.Get("created_from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid created_from")
			return
		}
		req.CreatedFrom = &t
	}

	if v := q.Get("created_to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid created_to")
			return
		}
		req.CreatedTo = &t
	}

	if v := q.Get("order"); v != "" {
		req.Order = v
	}

	result, err := service.GetAllItems(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't get items")
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

	var req dto.ConsumeItemRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	req.ID = id

	if err := service.ConsumeItem(r.Context(), req); err != nil {
		writeServiceError(w, err, "couldn't consume item")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*64)

	var req dto.CreateItemRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	item, err := service.CreateItem(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create item")
		return
	}

	response.WriteJSON(w, http.StatusOK, item)
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
