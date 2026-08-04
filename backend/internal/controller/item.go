package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := service.GetAllItems(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't get items")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	item, err := service.GetItem(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't get item")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func ConsumeItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024)

	var req dto.ConsumeItemRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
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
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	item, err := service.CreateItem(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create item")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
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
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.AddUpdateItemPropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.ItemID = itemId

	itemProp, err := service.AddItemProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't add item property")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(itemProp)
}

func UpdateItemProperty(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	propId, err := strconv.ParseInt(r.PathValue("prop_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*32)

	var req dto.AddUpdateItemPropertyRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.ItemID = itemId
	req.PropertyID = propId

	itemProp, err := service.UpdateItemProperty(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't update property of item")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(itemProp)
}

func RemoveItemProperty(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	propId, err := strconv.ParseInt(r.PathValue("prop_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	if err := service.RemoveItemProperty(r.Context(), itemId, propId); err != nil {
		writeServiceError(w, err, "couldn't remove item property")
		return
	}

	w.WriteHeader(http.StatusOK)
}
