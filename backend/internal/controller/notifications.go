package controller

import (
	"net/http"
	"strconv"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/pagination"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

// ListNotifications godoc
// @Summary List notifications
// @Tags Notifications
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Page size (default 20, max 50)"
// @Param offset query int false "Page offset (default 0)"
// @Success 200 {object} dto.NotificationsPage
// @Failure 400 {object} response.Error "invalid limit/offset"
// @Failure 401 {object} response.Error "not logged in"
// @Router /notifications [get]
func ListNotifications(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}

	limit, offset, err := pagination.GetLimitOffset(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	// Code pattern required for swaggo to not fail
	var notifications dto.NotificationsPage
	notifications, err = service.ListNotifications(r.Context(), id, limit, offset)
	if err != nil {
		writeServiceError(w, err, "couldn't list notifications")
		return
	}

	response.WriteJSON(w, http.StatusOK, notifications)
}

// ReadNotifications godoc
// @Summary Mark all notifications as read, returns amount of notifications read
// @Tags Notifications
// @Produce json
// @Security CookieAuth
// @Success 200 {object} int64
// @Failure 401 {object} response.Error "not logged in"
// @Router /notifications/read [post]
func ReadNotifications(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}

	count, err := service.ReadNotifications(r.Context(), id)
	if err != nil {
		writeServiceError(w, err, "couldn't read notifications")
		return
	}

	response.WriteJSON(w, http.StatusOK, count)
}

// DeleteNotification godoc
// @Summary Delete notification
// @Tags Notifications
// @Produce json
// @Security CookieAuth
// @Param id path int true "Notification ID"
// @Success 200
// @Failure 400 {object} response.Error "invalid notification id"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Failure 404 {object} response.Error "not found"
// @Router /notifications/{id} [delete]
func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "user ID not found in context")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	if err := service.DeleteNotification(r.Context(), id, userID); err != nil {
		writeServiceError(w, err, "couldn't delete notification")
		return
	}

	w.WriteHeader(http.StatusOK)
}
