package controller

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/response"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

// Left open, months is a cheap way to ask for an arbitrarily long axis and a scan of the whole
// audit log behind it.
var dashboardMonthRanges = []int32{3, 6, 12}

const defaultDashboardMonths = 6

// GetDashboard godoc
// @Summary Summarise expiry, consumption and stock (manager and admin only)
// @Description Every dashboard widget in one read. The window runs forward from the current month
// @Description for expiry, which asks what is about to go off, and backward from it for consumption,
// @Description which asks what has already been used. Both charts report a bucket per month over the
// @Description whole range, including the months in which nothing happened.
// @Tags Dashboard
// @Produce json
// @Security CookieAuth
// @Param months query int false "Range covered by both charts, one of 3, 6 or 12 (default 6)" Enums(3, 6, 12)
// @Param type_id query int false "Narrow every figure to one item type (default all types)"
// @Success 200 {object} dto.Dashboard
// @Failure 400 {object} response.Error "invalid query parameters"
// @Failure 401 {object} response.Error "not logged in"
// @Failure 403 {object} response.Error "forbidden"
// @Router /dashboard [get]
func GetDashboard(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	months := int32(defaultDashboardMonths)
	if val := query.Get("months"); val != "" {
		parsed, err := strconv.ParseInt(val, 10, 32)
		if err != nil || !slices.Contains(dashboardMonthRanges, int32(parsed)) {
			response.WriteError(w, http.StatusBadRequest, "invalid months")
			return
		}
		months = int32(parsed)
	}

	var typeID *int64
	if val := query.Get("type_id"); val != "" {
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid type_id")
			return
		}
		typeID = &parsed
	}

	// Code pattern required for swaggo to not fail
	var dashboard dto.Dashboard
	dashboard, err := service.GetDashboard(r.Context(), months, typeID)
	if err != nil {
		writeServiceError(w, err, "couldn't get the dashboard")
		return
	}

	response.WriteJSON(w, http.StatusOK, dashboard)
}
