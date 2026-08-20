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

// CreateItemRequest godoc
// @Summary Create an item request (admin only)
// @Tags ItemRequests
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body dto.ItemRequestCreateRequest true "Item request to create"
// @Success 200 {object} dto.ItemRequest
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /item-requests [post]
func CreateItemRequest(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*2)

	var req dto.ItemRequestCreateRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	itemRequest, err := service.CreateItemRequest(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create item request")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemRequest)
}

// CreatePersonalItemRequest godoc
// @Summary Create an item request on behalf of logged in user
// @Tags ItemRequests
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body dto.ItemRequestCreatePersonalRequest true "Item request to create"
// @Success 200 {object} dto.ItemRequest
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /item-requests/me [post]
func CreatePersonalItemRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}

	body := http.MaxBytesReader(w, r.Body, 1024*2)

	var personalReq dto.ItemRequestCreatePersonalRequest
	if err := json.NewDecoder(body).Decode(&personalReq); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	req := dto.ItemRequestCreateRequest{
		UserID: userID,
		ItemID: personalReq.ItemID,
		Reason: personalReq.Reason,
	}

	itemRequest, err := service.CreateItemRequest(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't create item request")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemRequest)
}

// GetPersonalItemRequest godoc
// @Summary Get personal item request (that the logged in user created) by item ID
// @Tags ItemRequests
// @Produce json
// @Security CookieAuth
// @Param item_id path int true "Item ID"
// @Success 200 {object} dto.ItemRequest
// @Failure 400 {object} response.Error "invalid item id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 404 {object} response.Error "not found"
// @Router /item-requests/me/{item_id} [get]
func GetPersonalItemRequest(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(r.PathValue("item_id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}

	itemRequest, err := service.GetItemRequest(r.Context(), dto.ItemRequestIdentifierRequest{
		ItemID: itemID,
		UserID: userID,
	})
	if err != nil {
		writeServiceError(w, err, "couldn't get item request")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemRequest)
}

// ApproveItemRequest godoc
// @Summary Approve an item request (admin only). Approving an item request deletes other unapproved requests for the same item.
// @Tags ItemRequests
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param body body dto.ItemRequestIdentifierRequest true "Item request to approve"
// @Success 200 {object} dto.ItemRequest
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /item-requests/approve [post]
func ApproveItemRequest(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*2)

	var req dto.ItemRequestIdentifierRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	itemRequest, err := service.ApproveItemRequest(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't approve item request")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemRequest)
}

// ListItemRequests godoc
// @Summary List item requests (admin only)
// @Tags ItemRequests
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Page size (default 20, max 50)"
// @Param offset query int false "Page offset (default 0)"
// @Param status query int false "Filter by item request status"
// @Success 200 {object} dto.ItemRequestsPage
// @Failure 400 {object} response.Error "invalid limit/offset"
// @Failure 401 {object} response.Error "not logged in"
// @Router /item-requests [get]
func ListItemRequests(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, offset, err := pagination.GetLimitOffset(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	listRequest := dto.ItemRequestsListRequest{
		Limit:  limit,
		Offset: offset,
	}

	if val := query.Get("status"); val != "" {
		listRequest.Status = &val
	}

	itemRequests, err := service.ListItemRequests(r.Context(), listRequest)
	if err != nil {
		writeServiceError(w, err, "couldn't list item requests")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemRequests)
}

// ListPersonalItemRequests godoc
// @Summary List personal item requests (of logged in user)
// @Tags ItemRequests
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Page size (default 20, max 50)"
// @Param offset query int false "Page offset (default 0)"
// @Success 200 {object} dto.ItemRequestsPage
// @Failure 400 {object} response.Error "invalid limit/offset"
// @Failure 401 {object} response.Error "not logged in"
// @Router /item-requests/me [get]
func ListPersonalItemRequests(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := pagination.GetLimitOffset(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}

	itemRequests, err := service.ListItemRequests(r.Context(), dto.ItemRequestsListRequest{
		Limit:  limit,
		Offset: offset,
		UserID: &userID,
	})
	if err != nil {
		writeServiceError(w, err, "couldn't list personal item requests")
		return
	}

	response.WriteJSON(w, http.StatusOK, itemRequests)
}

// DeleteItemRequest godoc
// @Summary Delete an item request (admin only)
// @Tags ItemRequests
// @Security CookieAuth
// @Param body body dto.ItemRequestIdentifierRequest true "Item request to delete"
// @Success 200
// @Failure 400 {object} response.Error "invalid request"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /item-requests [delete]
func DeleteItemRequest(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1024*2)

	var req dto.ItemRequestIdentifierRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := service.DeleteItemRequest(r.Context(), req); err != nil {
		writeServiceError(w, err, "couldn't delete request")
		return
	}

	w.WriteHeader(http.StatusOK)
}
