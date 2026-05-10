// Package api wires HTTP routes to handlers. Handlers are thin: they decode
// requests, call services, and encode responses. Business logic lives in
// services; SQL lives in repository.
package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/oroya/backend/internal/config"
	mw "github.com/oroya/backend/internal/middleware"
	"github.com/oroya/backend/internal/services"
)

type Deps struct {
	Cfg     *config.Config
	Logger  *slog.Logger
	Auth    *services.AuthService
	User    *services.UserService
	Video   *services.VideoService
	Upload  *services.UploadService
	Comment *services.CommentService
	Channel *services.ChannelService
	Search  *services.SearchService
	Admin   *services.AdminService
}

func NewRouter(d *Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(mw.RequestLogger(d.Logger))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   d.Cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Admin-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	rl := mw.NewRateLimiter(d.Cfg.RateLimit.RPS, d.Cfg.RateLimit.Burst)
	r.Use(rl.Middleware)

	requireAuth := mw.RequireAuth(d.Cfg.Supabase.JWTSecret)
	optionalAuth := mw.OptionalAuth(d.Cfg.Supabase.JWTSecret)
	requireAdmin := mw.RequireAdmin(d.Cfg.Admin.APIToken)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	infoH := NewInfoHandler(d.Cfg)
	authH := NewAuthHandler(d.Auth)
	uploadH := NewUploadHandler(d.Upload)
	videoH := NewVideoHandler(d.Video)
	commentH := NewCommentHandler(d.Comment)
	meH := NewMeHandler(d.User)
	channelH := NewChannelHandler(d.Channel)
	searchH := NewSearchHandler(d.Search)
	adminH := NewAdminHandler(d.Cfg, d.Admin)

	r.Route("/api/v1", func(r chi.Router) {
		// --- Frontend integration manifest ---
		r.Get("/info", infoH.Get)

		// --- Auth ---
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.Post("/refresh", authH.Refresh)
			r.Post("/google", authH.Google)
			r.Group(func(r chi.Router) {
				r.Use(requireAuth)
				r.Post("/logout", authH.Logout)
				r.Get("/me", authH.Me)
			})
		})

		// --- Me ---
		r.Group(func(r chi.Router) {
			r.Use(requireAuth)
			r.Get("/me", meH.Get)
			r.Put("/me", meH.Update)
		})

		// --- Videos ---
		r.Route("/videos", func(r chi.Router) {
			r.With(optionalAuth).Get("/", videoH.List)
			r.With(optionalAuth).Get("/{id}", videoH.Get)
			r.With(optionalAuth).Post("/{id}/view", videoH.View)

			r.Group(func(r chi.Router) {
				r.Use(requireAuth)
				r.Post("/upload-url", uploadH.SignURL)
				r.Post("/", videoH.Create)
				r.Put("/{id}", videoH.Update)
				r.Delete("/{id}", videoH.Delete)
				r.Post("/{id}/like", videoH.ToggleLike)

				r.Post("/{id}/links", videoH.AddLink)
				r.Delete("/{id}/links/{linkId}", videoH.DeleteLink)

				r.Post("/{id}/comments", commentH.Create)
			})
			r.Get("/{id}/links", videoH.ListLinks)
			r.Get("/{id}/comments", commentH.List)
		})

		// --- Comments ---
		r.Group(func(r chi.Router) {
			r.Use(requireAuth)
			r.Delete("/comments/{id}", commentH.Delete)
			r.Post("/comments/{id}/like", commentH.ToggleLike)
		})

		// --- Channels ---
		r.With(optionalAuth).Get("/channels/{id}", channelH.Get)
		r.With(requireAuth).Post("/channels/{id}/subscribe", channelH.ToggleSubscribe)

		// --- Search ---
		r.Get("/search", searchH.Search)

		// --- Admin ---
		r.Route("/admin", func(r chi.Router) {
			r.Get("/health", adminH.Health) // public
			r.Group(func(r chi.Router) {
				r.Use(requireAdmin)
				r.Get("/stats", adminH.Stats)
				r.Get("/queue", adminH.Queue)
				r.Get("/worker-status", adminH.WorkerStatus)
				r.Get("/storage-status", adminH.StorageStatus)
			})
		})
	})

	r.Get("/admin", adminH.Dashboard)

	return r
}
