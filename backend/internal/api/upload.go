package api

import (
	"net/http"

	mw "github.com/oroya/backend/internal/middleware"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/utils"
)

type UploadHandler struct {
	svc *services.UploadService
}

func NewUploadHandler(svc *services.UploadService) *UploadHandler { return &UploadHandler{svc: svc} }

type uploadURLRequest struct {
	Filename string `json:"filename"`
}

func (h *UploadHandler) SignURL(w http.ResponseWriter, r *http.Request) {
	claims, _ := mw.ClaimsFrom(r.Context())
	var req uploadURLRequest
	_ = utils.DecodeJSON(r, &req)
	resp, err := h.svc.SignUploadURL(claims.UserID, req.Filename)
	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "sign_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}
