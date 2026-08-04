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
		http.Error(w, "couldn't get items", http.StatusInternalServerError)
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
		http.Error(w, "couldn't get item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	item, err := service.CreateItem(r.Context())
	if err != nil {
		http.Error(w, "couldn't create item", http.StatusInternalServerError)
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
		http.Error(w, "couldn't delete item", http.StatusInternalServerError)
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
		http.Error(w, "couldn't add item property", http.StatusInternalServerError)
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
		http.Error(w, "couldn't update property of item", http.StatusInternalServerError)
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
		http.Error(w, "couldn't remove item property", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
