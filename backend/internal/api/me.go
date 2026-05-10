package api

import (
	"net/http"

	mw "github.com/oroya/backend/internal/middleware"
	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/utils"
)

type MeHandler struct {
	svc *services.UserService
}

func NewMeHandler(svc *services.UserService) *MeHandler { return &MeHandler{svc: svc} }

func (h *MeHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	p, err := h.svc.Get(r.Context(), claims.UserID)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, p)
}

func (h *MeHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	var req models.UpdateProfileRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	p, err := h.svc.Update(r.Context(), claims.UserID, &req)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, p)
}
