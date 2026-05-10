package api

import (
	"errors"
	"net/http"
	"strings"

	mw "github.com/oroya/backend/internal/middleware"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/supabase"
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
		var sbErr *supabase.APIError
		if errors.As(err, &sbErr) && strings.Contains(strings.ToLower(sbErr.Raw), "row-level security") {
			utils.WriteError(w, http.StatusBadGateway, "storage_policy_denied", "storage bucket policy is blocking uploads")
			return
		}
		utils.WriteError(w, http.StatusBadGateway, "sign_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}
