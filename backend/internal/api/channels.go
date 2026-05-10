package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	mw "github.com/oroya/backend/internal/middleware"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/utils"
)

type ChannelHandler struct {
	svc *services.ChannelService
}

func NewChannelHandler(svc *services.ChannelService) *ChannelHandler {
	return &ChannelHandler{svc: svc}
}

func (h *ChannelHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	viewerID := ""
	if c, ok := mw.ClaimsFrom(r.Context()); ok {
		viewerID = c.UserID
	}
	view, err := h.svc.Get(r.Context(), id, viewerID)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, view)
}

func (h *ChannelHandler) ToggleSubscribe(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	id := chi.URLParam(r, "id")
	subscribed, err := h.svc.ToggleSubscribe(r.Context(), claims.UserID, id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrSelfSubscribe):
			utils.WriteError(w, http.StatusBadRequest, "self_subscribe", err.Error())
		default:
			writeRepoError(w, err)
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"channel_id": id,
		"subscribed": subscribed,
	})
}
