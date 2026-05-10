package api

import (
	"errors"
	"net/http"
	"strings"

	mw "github.com/oroya/backend/internal/middleware"
	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/supabase"
	"github.com/oroya/backend/internal/utils"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	resp, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	resp, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	resp, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	tok := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if err := h.svc.Logout(r.Context(), strings.TrimSpace(tok)); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "logout_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	profile, err := h.svc.Me(r.Context(), claims.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "not_found", "profile not found")
		return
	}
	utils.WriteJSON(w, http.StatusOK, profile)
}

func (h *AuthHandler) Google(w http.ResponseWriter, r *http.Request) {
	// Placeholder: full Google OAuth flow will be implemented when frontend integration lands.
	utils.WriteError(w, http.StatusNotImplemented, "not_implemented", "google oauth pending")
}

func writeAuthError(w http.ResponseWriter, err error) {
	var sbErr *supabase.APIError
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		utils.WriteError(w, http.StatusBadRequest, "invalid_input", err.Error())
	case errors.Is(err, services.ErrUsernameTaken):
		utils.WriteError(w, http.StatusConflict, "username_taken", err.Error())
	case errors.As(err, &sbErr) && sbErr.Status >= 500:
		utils.WriteError(w, http.StatusBadGateway, "auth_service_failed", "authentication service could not create the user")
	case errors.As(err, &sbErr) && sbErr.Status == http.StatusUnauthorized:
		utils.WriteError(w, http.StatusBadGateway, "auth_service_unauthorized", "authentication service rejected the configured API key")
	default:
		utils.WriteError(w, http.StatusUnauthorized, "auth_failed", err.Error())
	}
}
