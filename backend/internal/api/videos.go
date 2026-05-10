package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	mw "github.com/oroya/backend/internal/middleware"
	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/repository"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/utils"
)

type VideoHandler struct {
	svc *services.VideoService
}

func NewVideoHandler(svc *services.VideoService) *VideoHandler { return &VideoHandler{svc: svc} }

func (h *VideoHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	vids, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"videos": vids})
}

func (h *VideoHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	vids, err := h.svc.ListByOwner(r.Context(), claims.UserID, limit, offset)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"videos": vids})
}

func (h *VideoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	v, assets, links, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"video":  v,
		"assets": assets,
		"links":  links,
	})
}

func (h *VideoHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	var req models.CreateVideoRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	v, err := h.svc.Create(r.Context(), claims.UserID, &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, v)
}

func (h *VideoHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	var req models.UpdateVideoRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	v, err := h.svc.Update(r.Context(), id, claims.UserID, &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, v)
}

func (h *VideoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id, claims.UserID); err != nil {
		writeServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *VideoHandler) View(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var userID *string
	if claims, ok := mw.ClaimsFrom(r.Context()); ok {
		userID = &claims.UserID
	}
	ipHash := utils.HashIP(r.RemoteAddr, r.Header.Get("X-Forwarded-For"))
	if err := h.svc.RegisterView(r.Context(), id, userID, ipHash); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "view_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "ok"})
}

func (h *VideoHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	liked, err := h.svc.ToggleLike(r.Context(), id, claims.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "like_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}

func (h *VideoHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	links, err := h.svc.ListLinks(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "links_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"links": links})
}

func (h *VideoHandler) AddLink(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	var req models.CreateVideoLinkRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	link, err := h.svc.AddLink(r.Context(), id, claims.UserID, &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, link)
}

func (h *VideoHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	linkID := chi.URLParam(r, "linkId")
	if err := h.svc.DeleteLink(r.Context(), id, linkID, claims.UserID); err != nil {
		writeServiceError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func writeRepoError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, "not_found", "resource not found")
		return
	}
	utils.WriteError(w, http.StatusInternalServerError, "internal", err.Error())
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		utils.WriteError(w, http.StatusBadRequest, "invalid_input", err.Error())
	case errors.Is(err, services.ErrForbidden):
		utils.WriteError(w, http.StatusForbidden, "forbidden", err.Error())
	case errors.Is(err, repository.ErrNotFound):
		utils.WriteError(w, http.StatusNotFound, "not_found", err.Error())
	default:
		utils.WriteError(w, http.StatusInternalServerError, "internal", err.Error())
	}
}
