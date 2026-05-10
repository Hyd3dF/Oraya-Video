package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	mw "github.com/oroya/backend/internal/middleware"
	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/utils"
)

type CommentHandler struct {
	svc *services.CommentService
}

func NewCommentHandler(svc *services.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	cs, err := h.svc.List(r.Context(), id, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"comments": cs})
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	var req models.CreateCommentRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	c, err := h.svc.Create(r.Context(), id, claims.UserID, &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, c)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id, claims.UserID); err != nil {
		writeServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *CommentHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	liked, err := h.svc.ToggleLike(r.Context(), id, claims.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "like_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}
