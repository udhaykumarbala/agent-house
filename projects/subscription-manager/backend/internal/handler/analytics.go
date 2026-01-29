package handler

import (
	"net/http"
	"strconv"

	"github.com/subtrack/backend/internal/service"
	"github.com/subtrack/backend/pkg/response"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	summary, err := h.analyticsService.GetSummary(r.Context(), user.ID)
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, summary)
}

func (h *AnalyticsHandler) Upcoming(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	groups, err := h.analyticsService.GetUpcoming(r.Context(), user.ID, days)
	if err != nil {
		response.InternalError(w)
		return
	}

	totalUpcoming := 0.0
	for _, g := range groups {
		totalUpcoming += g.Total
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"groups":   groups,
		"total":    totalUpcoming,
		"daysSpan": days,
	})
}
