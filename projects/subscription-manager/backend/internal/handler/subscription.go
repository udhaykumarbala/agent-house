package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/subtrack/backend/internal/model"
	"github.com/subtrack/backend/internal/service"
	"github.com/subtrack/backend/pkg/response"
)

type SubscriptionHandler struct {
	subService *service.SubscriptionService
}

func NewSubscriptionHandler(subService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subService: subService}
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	subs, err := h.subService.GetAll(r.Context(), user.ID)
	if err != nil {
		response.InternalError(w)
		return
	}

	if subs == nil {
		subs = make([]*model.Subscription, 0)
	}

	response.JSONWithMeta(w, http.StatusOK, subs, &response.Meta{Total: len(subs)})
}

func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid subscription ID")
		return
	}

	sub, err := h.subService.GetByID(r.Context(), id, user.ID)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, sub)
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	var input model.SubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	sub, err := h.subService.Create(r.Context(), user.ID, &input)
	if err != nil {
		handleValidationError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, sub)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid subscription ID")
		return
	}

	var input model.SubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	sub, err := h.subService.Update(r.Context(), id, user.ID, &input)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			response.NotFound(w)
			return
		}
		handleValidationError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, sub)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid subscription ID")
		return
	}

	err = h.subService.Delete(r.Context(), id, user.ID)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Subscription deleted"})
}

func handleValidationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidName):
		response.ValidationError(w, []response.FieldError{{Field: "name", Message: "Name must be 1-100 characters"}})
	case errors.Is(err, service.ErrInvalidPrice):
		response.ValidationError(w, []response.FieldError{{Field: "price", Message: "Price must be between 0 and 999999.99"}})
	case errors.Is(err, service.ErrInvalidCycle):
		response.ValidationError(w, []response.FieldError{{Field: "billingCycle", Message: "Invalid billing cycle"}})
	case errors.Is(err, service.ErrInvalidCategory):
		response.ValidationError(w, []response.FieldError{{Field: "category", Message: "Invalid category"}})
	case errors.Is(err, service.ErrInvalidDate):
		response.ValidationError(w, []response.FieldError{{Field: "nextBillingDate", Message: "Invalid date format (use YYYY-MM-DD)"}})
	case errors.Is(err, service.ErrInvalidColor):
		response.ValidationError(w, []response.FieldError{{Field: "color", Message: "Invalid color format (use #RRGGBB)"}})
	case errors.Is(err, service.ErrInvalidReminder):
		response.ValidationError(w, []response.FieldError{{Field: "reminderDays", Message: "Reminder days must be 0-30"}})
	default:
		response.InternalError(w)
	}
}
