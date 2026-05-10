package api

import (
	"net/http"
	"runtime"
	"time"

	"github.com/oroya/backend/internal/config"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/utils"
)

type AdminHandler struct {
	cfg     *config.Config
	svc     *services.AdminService
	started time.Time
}

func NewAdminHandler(cfg *config.Config, svc *services.AdminService) *AdminHandler {
	return &AdminHandler{cfg: cfg, svc: svc, started: time.Now()}
}

func (h *AdminHandler) Health(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"status":    "ok",
		"uptime_s":  int(time.Since(h.started).Seconds()),
		"timestamp": time.Now().UTC(),
	}
	if h.svc != nil {
		health := h.svc.Health(r.Context())
		resp["status"] = health.Status
		resp["database"] = health.Database
		resp["storage"] = health.Storage
		resp["worker"] = health.Worker
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	resp := map[string]any{
		"runtime": map[string]any{
			"goroutines":     runtime.NumGoroutine(),
			"memory_alloc_b": m.Alloc,
			"memory_sys_b":   m.Sys,
			"gc_cycles":      m.NumGC,
			"uptime_s":       int(time.Since(h.started).Seconds()),
		},
	}
	if stats, err := h.svc.Stats(r.Context()); err == nil {
		resp["data"] = stats
	} else {
		resp["data_error"] = err.Error()
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *AdminHandler) Queue(w http.ResponseWriter, r *http.Request) {
	q, err := h.svc.Queue(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "queue_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, q)
}

func (h *AdminHandler) WorkerStatus(w http.ResponseWriter, r *http.Request) {
	st := h.svc.WorkerStatus()
	if st == nil {
		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"running":     false,
			"concurrency": h.cfg.Worker.Concurrency,
			"ffmpeg_bin":  h.cfg.Worker.FFmpegBin,
			"temp_dir":    h.cfg.Worker.TempDir,
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"running":     true,
		"concurrency": st.Concurrency,
		"pending":     st.Pending,
		"processing":  st.Processing,
		"completed":   st.Completed,
		"failed":      st.Failed,
		"current":     st.Current,
		"ffmpeg_bin":  h.cfg.Worker.FFmpegBin,
		"temp_dir":    h.cfg.Worker.TempDir,
	})
}

func (h *AdminHandler) StorageStatus(w http.ResponseWriter, r *http.Request) {
	st, err := h.svc.StorageStatus(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "storage_failed", err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, st)
}


func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardHTML))
}
