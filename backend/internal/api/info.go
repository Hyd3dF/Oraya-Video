package api

import (
	"net/http"

	"github.com/oroya/backend/internal/config"
	"github.com/oroya/backend/internal/utils"
)

// InfoHandler exposes a single, frontend-facing manifest at /api/v1/info.
// It is the contract: anyone integrating with this backend reads this and
// nothing else (the frontend never sees Supabase keys, DB URLs, or secrets).
type InfoHandler struct {
	cfg *config.Config
}

func NewInfoHandler(cfg *config.Config) *InfoHandler { return &InfoHandler{cfg: cfg} }

func (h *InfoHandler) Get(w http.ResponseWriter, r *http.Request) {
	base := h.cfg.Server.PublicURL
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"name":         "oroya-backend",
		"api_base_url": base,
		"frontend_env": map[string]string{
			"PUBLIC_API_BASE_URL": base,
		},
		"cors_allowed_origins": h.cfg.CORS.AllowedOrigins,
		"auth_token": map[string]string{
			"header": "Authorization",
			"format": "Bearer <access_token>",
			"obtain": "POST /api/v1/auth/login or /api/v1/auth/register — response field 'access_token'",
			"refresh": "POST /api/v1/auth/refresh with {\"refresh_token\": \"...\"} when access_token expires",
		},
		"endpoints": endpointMap(),
		"upload_flow": []string{
			"1) POST /api/v1/videos/upload-url   -> returns { upload_url, storage_path }",
			"2) PUT  <upload_url>  with raw video bytes (frontend uploads directly to Supabase Storage)",
			"3) POST /api/v1/videos               -> { title, description, visibility, storage_path } -> returns video.id with status='processing'",
			"4) Poll GET /api/v1/videos/{id} until status='ready'; assets[].playlist_url is the HLS URL for HLS.js",
		},
	})
}

func endpointMap() map[string]any {
	return map[string]any{
		"auth": map[string]any{
			"register": "POST /api/v1/auth/register   { email, password, real_name, username }",
			"login":    "POST /api/v1/auth/login      { email, password }",
			"refresh":  "POST /api/v1/auth/refresh    { refresh_token }",
			"logout":   "POST /api/v1/auth/logout     (auth required)",
			"me":       "GET  /api/v1/auth/me         (auth required)",
		},
		"profile": map[string]any{
			"get_me":    "GET  /api/v1/me              (auth required)",
			"update_me": "PUT  /api/v1/me              (auth required) { real_name?, username?, display_name?, avatar_url?, banner_url?, bio? }",
		},
		"videos": map[string]any{
			"list":         "GET    /api/v1/videos?limit=24&offset=0",
			"get":          "GET    /api/v1/videos/{id}",
			"create":       "POST   /api/v1/videos                 (auth) { title, description, visibility, storage_path, duration_seconds }",
			"update":       "PUT    /api/v1/videos/{id}            (auth, owner only) { title?, description?, visibility? }",
			"delete":       "DELETE /api/v1/videos/{id}            (auth, owner only)",
			"view":         "POST   /api/v1/videos/{id}/view",
			"like_toggle":  "POST   /api/v1/videos/{id}/like       (auth) -> { liked: bool }",
			"upload_url":   "POST   /api/v1/videos/upload-url      (auth) { filename } -> { upload_url, storage_path, expires_at }",
			"links_list":   "GET    /api/v1/videos/{id}/links",
			"links_add":    "POST   /api/v1/videos/{id}/links      (auth, owner) { title, url, sort_order }",
			"links_delete": "DELETE /api/v1/videos/{id}/links/{linkId} (auth, owner)",
		},
		"comments": map[string]any{
			"list":        "GET    /api/v1/videos/{id}/comments?limit=30&offset=0",
			"create":      "POST   /api/v1/videos/{id}/comments   (auth) { content, parent_id? }",
			"delete":      "DELETE /api/v1/comments/{id}          (auth, owner only)",
			"like_toggle": "POST   /api/v1/comments/{id}/like     (auth) -> { liked: bool }",
		},
		"channels": map[string]any{
			"get":              "GET  /api/v1/channels/{id}",
			"toggle_subscribe": "POST /api/v1/channels/{id}/subscribe (auth) -> { subscribed: bool }",
		},
		"search": map[string]any{
			"all": "GET /api/v1/search?q=<query>&limit=25 -> { videos: [...], channels: [...] }",
		},
		"system": map[string]any{
			"info":   "GET /api/v1/info        (this endpoint)",
			"health": "GET /api/v1/admin/health (public)",
		},
	}
}
