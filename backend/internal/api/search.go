package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/utils"
)

type SearchHandler struct {
	svc *services.SearchService
}

func NewSearchHandler(svc *services.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	results, err := h.svc.Search(r.Context(), q, limit)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			utils.WriteError(w, http.StatusBadRequest, "missing_query", "q parameter required")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "search_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, results)
}
