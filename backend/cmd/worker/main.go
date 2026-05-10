// Command oroya-worker runs the background video processing pipeline as a
// standalone process. With the MVP in-memory queue, server and worker share
// a process — so this binary is only useful once the queue is swapped for
// Redis (or another shared backend).
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/oroya/backend/internal/config"
	"github.com/oroya/backend/internal/repository"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sb := supabase.New(cfg.Supabase.URL, cfg.Supabase.ServiceRoleKey, cfg.Supabase.AnonKey)
	if err := sb.PingREST(ctx); err != nil {
		logger.Error("supabase REST ping failed", "err", err)
		os.Exit(1)
	}

	db := repository.New(sb)
	defer db.Close()
	videos := repository.NewVideoRepository(db)

	queue := worker.NewMemoryQueue() // replace with Redis-backed queue for HA
	ff := worker.NewFFmpeg(cfg.Worker.FFmpegBin, cfg.Worker.FFprobeBin)
	processor := worker.NewProcessor(
		queue, videos, ff, sb,
		cfg.Storage.BucketRaw, cfg.Storage.BucketHLS,
		cfg.Worker.TempDir, cfg.Worker.Concurrency, logger,
	)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		logger.Info("worker shutdown initiated")
		cancel()
	}()

	logger.Info("oroya-worker started", "concurrency", cfg.Worker.Concurrency)
	processor.Run(ctx)
	logger.Info("oroya-worker stopped")
}
