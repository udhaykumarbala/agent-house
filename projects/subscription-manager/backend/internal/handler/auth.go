package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/subtrack/backend/internal/model"
	"github.com/subtrack/backend/internal/service"
	"github.com/subtrack/backend/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input model.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	user, token, err := h.authService.Register(r.Context(), &input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmail):
			response.ValidationError(w, []response.FieldError{{Field: "email", Message: "Invalid email address"}})
		case errors.Is(err, service.ErrPasswordTooShort):
			response.ValidationError(w, []response.FieldError{{Field: "password", Message: "Password must be at least 12 characters"}})
		case errors.Is(err, service.ErrPasswordTooLong):
			response.ValidationError(w, []response.FieldError{{Field: "password", Message: "Password must be at most 128 characters"}})
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid email or password")
		default:
			response.InternalError(w)
		}
		return
	}

	setSessionCookie(w, token)
	response.JSON(w, http.StatusCreated, user.ToResponse())
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	user, token, err := h.authService.Login(r.Context(), &input, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid email or password")
			return
		}
		response.InternalError(w)
		return
	}

	setSessionCookie(w, token)
	response.JSON(w, http.StatusOK, user.ToResponse())
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		h.authService.Logout(r.Context(), cookie.Value)
	}

	clearSessionCookie(w)
	response.JSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	response.JSON(w, http.StatusOK, user.ToResponse())
}

func (h *AuthHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	var prefs model.UserPreferences
	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if err := h.authService.UpdatePreferences(r.Context(), user.ID, &prefs); err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, &prefs)
}

func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if user == nil {
		response.Unauthorized(w)
		return
	}

	if err := h.authService.DeleteAccount(r.Context(), user.ID); err != nil {
		response.InternalError(w)
		return
	}

	clearSessionCookie(w)
	response.JSON(w, http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 60 * 60,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
