// Command oroya-server runs the HTTP API and (for MVP) the in-process video
// processing worker. The backend talks to Supabase exclusively via REST:
// Auth (/auth/v1/), database (/rest/v1/ — PostgREST), Storage (/storage/v1/).
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oroya/backend/internal/api"
	"github.com/oroya/backend/internal/config"
	"github.com/oroya/backend/internal/repository"
	"github.com/oroya/backend/internal/services"
	"github.com/oroya/backend/internal/supabase"
	"github.com/oroya/backend/internal/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "err", err)
		os.Exit(1)
	}

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	// --- Supabase client (sole DB/Auth/Storage gateway) ---
	sb := supabase.New(cfg.Supabase.URL, cfg.Supabase.ServiceRoleKey, cfg.Supabase.AnonKey)

	pingCtx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
	if err := sb.PingREST(pingCtx); err != nil {
		cancel()
		logger.Error("supabase REST ping failed", "err", err)
		os.Exit(1)
	}
	cancel()
	logger.Info("supabase reachable", "url", cfg.Supabase.URL)

	db := repository.New(sb)
	defer db.Close()

	// --- Repositories ---
	users := repository.NewUserRepository(db)
	videos := repository.NewVideoRepository(db)
	comments := repository.NewCommentRepository(db)
	subs := repository.NewSubscriptionRepository(db)
	views := repository.NewViewRepository(db)
	search := repository.NewSearchRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// --- Worker ---
	queue := worker.NewMemoryQueue()
	ff := worker.NewFFmpeg(cfg.Worker.FFmpegBin, cfg.Worker.FFprobeBin)
	processVideos := true
	if err := ff.Available(); err != nil {
		processVideos = false
		logger.Warn("ffmpeg unavailable; uploaded source files will be served directly", "err", err)
	}
	processor := worker.NewProcessor(
		queue, videos, ff, sb,
		cfg.Storage.BucketRaw, cfg.Storage.BucketHLS,
		cfg.Worker.TempDir, cfg.Worker.Concurrency, logger,
	)
	go processor.Run(rootCtx)
	logger.Info("worker started", "concurrency", cfg.Worker.Concurrency)

	// --- Services ---
	authSvc := services.NewAuthService(users, sb)
	userSvc := services.NewUserService(users)
	uploadSvc := services.NewUploadService(sb, cfg.Storage.BucketRaw)
	videoSvc := services.NewVideoService(videos, views, queue, sb, cfg.Storage.BucketRaw, processVideos)
	commentSvc := services.NewCommentService(comments)
	channelSvc := services.NewChannelService(users, videos, subs)
	searchSvc := services.NewSearchService(search)
	adminSvc := services.NewAdminService(adminRepo, processor, sb, services.AdminBuckets{
		Raw:        cfg.Storage.BucketRaw,
		HLS:        cfg.Storage.BucketHLS,
		Thumbnails: cfg.Storage.BucketThumbnails,
	}, db)

	router := api.NewRouter(&api.Deps{
		Cfg:     cfg,
		Logger:  logger,
		Auth:    authSvc,
		User:    userSvc,
		Video:   videoSvc,
		Upload:  uploadSvc,
		Comment: commentSvc,
		Channel: channelSvc,
		Search:  searchSvc,
		Admin:   adminSvc,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		logger.Info("oroya-server listening", "addr", srv.Addr, "env", cfg.Server.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "err", err)
			rootCancel()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("shutdown initiated")

	shutdownCtx, scancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer scancel()
	_ = srv.Shutdown(shutdownCtx)
	rootCancel()
	queue.Close()
	logger.Info("shutdown complete")
}
