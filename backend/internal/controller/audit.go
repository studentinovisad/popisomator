package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/pagination"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

// ListAuditLog godoc
// @Summary List recorded changes (admin only)
// @Tags AuditLog
// @Produce json
// @Security CookieAuth
// @Param limit query int false "Page size (default 20, max 50)"
// @Param offset query int false "Page offset (default 0)"
// @Param action query string false "Filter by action" Enums(item_create, item_update, item_consume, item_delete, item_property_add, item_property_update, item_property_remove, item_request_create, item_request_approve, item_request_delete, item_request_supersede, item_type_create, item_type_update, item_type_delete, item_type_property_add, item_type_property_update, item_type_property_remove, item_type_property_reorder, property_create, property_update, property_delete)
// @Param target_type query string false "Filter by the kind of entity changed" Enums(item, item_type, property)
// @Param target_id query int false "Filter to one entity's own history; requires target_type"
// @Param actor_id query int false "Filter by who made the change"
// @Param created_from query string false "Filter by when the change was made, RFC3339"
// @Param created_to query string false "Filter by when the change was made, RFC3339"
// @Success 200 {object} dto.AuditLogPage
// @Failure 400 {object} response.Error "invalid query parameters"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /audit-log [get]
func ListAuditLog(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, offset, err := pagination.GetLimitOffset(r)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid limit/offset")
		return
	}

	req := dto.ListAuditLogRequest{Limit: limit, Offset: offset}

	// The action and target type are Postgres enums. Rejecting an unknown one here keeps it from
	// reaching the database, where the cast would fail as a 500 rather than a bad request.
	if val := query.Get("action"); val != "" {
		action, ok := dto.ParseAuditAction(val)
		if !ok {
			response.WriteError(w, http.StatusBadRequest, "invalid action")
			return
		}
		req.Action = &action
	}

	if val := query.Get("target_type"); val != "" {
		targetType, ok := dto.ParseAuditTargetType(val)
		if !ok {
			response.WriteError(w, http.StatusBadRequest, "invalid target_type")
			return
		}
		req.TargetType = &targetType
	}

	if val := query.Get("target_id"); val != "" {
		targetID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid target_id")
			return
		}
		// Ids are only unique within a kind of entity, so an id on its own would mix an item's
		// history with the item type that happens to share its number.
		if req.TargetType == nil {
			response.WriteError(w, http.StatusBadRequest, "target_id requires target_type")
			return
		}
		req.TargetID = &targetID
	}

	if val := query.Get("actor_id"); val != "" {
		actorID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid actor_id")
			return
		}
		req.ActorID = &actorID
	}

	if val := query.Get("created_from"); val != "" {
		parsedTime, err := time.Parse(time.RFC3339, val)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid created_from")
			return
		}
		req.CreatedFrom = &parsedTime
	}

	if val := query.Get("created_to"); val != "" {
		parsedTime, err := time.Parse(time.RFC3339, val)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid created_to")
			return
		}
		req.CreatedTo = &parsedTime
	}

	// Code pattern required for swaggo to not fail
	var auditLog dto.AuditLogPage
	auditLog, err = service.ListAuditLog(r.Context(), req)
	if err != nil {
		writeServiceError(w, err, "couldn't list the audit log")
		return
	}

	response.WriteJSON(w, http.StatusOK, auditLog)
}

// ListAuditLogActors godoc
// @Summary List everyone who has made a recorded change (admin only)
// @Tags AuditLog
// @Produce json
// @Security CookieAuth
// @Success 200 {array} dto.AuditActorOption
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /audit-log/actors [get]
func ListAuditLogActors(w http.ResponseWriter, r *http.Request) {
	actors, err := service.ListAuditLogActors(r.Context())
	if err != nil {
		writeServiceError(w, err, "couldn't list audit log actors")
		return
	}

	response.WriteJSON(w, http.StatusOK, actors)
}
