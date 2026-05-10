package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/repository"
	"github.com/oroya/backend/internal/supabase"
)

// Processor consumes Jobs and runs the full pipeline.
type Processor struct {
	queue       Queue
	videos      repository.VideoRepository
	ffmpeg      *FFmpeg
	sb          *supabase.Client
	rawBucket   string
	hlsBucket   string
	tempDir     string
	concurrency int
	logger      *slog.Logger

	stats workerStats
}

func NewProcessor(
	q Queue,
	videos repository.VideoRepository,
	ff *FFmpeg,
	sb *supabase.Client,
	rawBucket, hlsBucket, tempDir string,
	concurrency int,
	logger *slog.Logger,
) *Processor {
	return &Processor{
		queue:       q,
		videos:      videos,
		ffmpeg:      ff,
		sb:          sb,
		rawBucket:   rawBucket,
		hlsBucket:   hlsBucket,
		tempDir:     tempDir,
		concurrency: concurrency,
		logger:      logger,
	}
}

type workerStats struct {
	mu         sync.RWMutex
	processing int64
	completed  int64
	failed     int64
	current    map[string]string // "worker-N" -> video id
}

type Status struct {
	Concurrency int               `json:"concurrency"`
	Pending     int               `json:"pending"`
	Processing  int64             `json:"processing"`
	Completed   int64             `json:"completed"`
	Failed      int64             `json:"failed"`
	Current     map[string]string `json:"current"`
}

func (p *Processor) Status() Status {
	p.stats.mu.RLock()
	defer p.stats.mu.RUnlock()
	cur := make(map[string]string, len(p.stats.current))
	for k, v := range p.stats.current {
		cur[k] = v
	}
	return Status{
		Concurrency: p.concurrency,
		Pending:     p.queue.Len(),
		Processing:  p.stats.processing,
		Completed:   p.stats.completed,
		Failed:      p.stats.failed,
		Current:     cur,
	}
}

func (p *Processor) Run(ctx context.Context) {
	if err := os.MkdirAll(p.tempDir, 0o755); err != nil {
		p.logger.Error("worker: mkdir tempdir failed", "err", err)
		return
	}

	p.stats.current = make(map[string]string)

	var wg sync.WaitGroup
	for i := 0; i < p.concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p.workerLoop(ctx, id)
		}(i)
	}

	<-ctx.Done()
	p.queue.Close()
	wg.Wait()
}

func (p *Processor) workerLoop(ctx context.Context, id int) {
	for {
		job, ok := p.queue.Dequeue()
		if !ok {
			return
		}
		p.markStart(id, job.VideoID)
		p.logger.Info("worker: processing", "worker", id, "video_id", job.VideoID)

		err := p.process(ctx, job)
		p.markEnd(id, err == nil)

		if err != nil {
			p.logger.Error("worker: job failed", "video_id", job.VideoID, "err", err)
			if setErr := p.videos.SetStatus(ctx, job.VideoID, "failed"); setErr != nil {
				p.logger.Error("worker: set failed status failed", "err", setErr)
			}
			continue
		}
		if err := p.videos.SetStatus(ctx, job.VideoID, "ready"); err != nil {
			p.logger.Error("worker: set ready status failed", "err", err)
		}
	}
}

func (p *Processor) markStart(id int, videoID string) {
	p.stats.mu.Lock()
	p.stats.processing++
	p.stats.current[workerKey(id)] = videoID
	p.stats.mu.Unlock()
}

func (p *Processor) markEnd(id int, ok bool) {
	p.stats.mu.Lock()
	p.stats.processing--
	delete(p.stats.current, workerKey(id))
	if ok {
		p.stats.completed++
	} else {
		p.stats.failed++
	}
	p.stats.mu.Unlock()
}

func (p *Processor) process(ctx context.Context, job Job) error {
	if p.sb == nil {
		return errors.New("storage client not configured")
	}

	workDir := filepath.Join(p.tempDir, job.VideoID)
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return err
	}
	defer os.RemoveAll(workDir)

	// 1. Download raw video from raw-videos bucket.
	ext := filepath.Ext(job.StoragePath)
	if ext == "" {
		ext = ".mp4"
	}
	rawPath := filepath.Join(workDir, "raw"+ext)

	dlCtx, cancelDL := context.WithTimeout(ctx, 30*time.Minute)
	if err := p.sb.Download(dlCtx, p.rawBucket, job.StoragePath, rawPath); err != nil {
		cancelDL()
		return fmt.Errorf("download raw: %w", err)
	}
	cancelDL()

	// 2. Transcode to HLS (720p MVP).
	hlsDir := filepath.Join(workDir, "hls")
	if err := os.MkdirAll(hlsDir, 0o755); err != nil {
		return err
	}
	transcodeCtx, cancelTC := context.WithTimeout(ctx, 60*time.Minute)
	if err := p.ffmpeg.TranscodeHLS(transcodeCtx, rawPath, hlsDir, MVPRendition720p); err != nil {
		cancelTC()
		return fmt.Errorf("transcode: %w", err)
	}
	cancelTC()

	// 3. Upload all .m3u8 + .ts files into hls-videos/{video_id}/720p/
	prefix := fmt.Sprintf("%s/%s", job.VideoID, MVPRendition720p.Quality)
	totalSize, err := p.uploadDir(ctx, hlsDir, prefix)
	if err != nil {
		return fmt.Errorf("upload hls: %w", err)
	}

	// 4. Register asset row.
	asset := &models.VideoAsset{
		ID:          uuid.NewString(),
		VideoID:     job.VideoID,
		Quality:     MVPRendition720p.Quality,
		PlaylistURL: p.sb.PublicURL(p.hlsBucket, prefix+"/playlist.m3u8"),
		MasterURL:   p.sb.PublicURL(p.hlsBucket, prefix+"/master.m3u8"),
		Width:       MVPRendition720p.Width,
		Height:      MVPRendition720p.Height,
		Bitrate:     MVPRendition720p.Bitrate,
		SizeBytes:   totalSize,
		CreatedAt:   time.Now().UTC(),
	}
	if err := p.videos.AddAsset(ctx, asset); err != nil {
		return fmt.Errorf("add asset: %w", err)
	}

	// 5. Best-effort cleanup of the raw upload (HLS is the source of truth now).
	rmCtx, cancelRM := context.WithTimeout(ctx, 30*time.Second)
	if err := p.sb.Remove(rmCtx, p.rawBucket, job.StoragePath); err != nil {
		p.logger.Warn("worker: raw cleanup failed", "video_id", job.VideoID, "err", err)
	}
	cancelRM()

	return nil
}

func workerKey(id int) string { return "worker-" + strconv.Itoa(id) }

func (p *Processor) uploadDir(ctx context.Context, srcDir, prefix string) (int64, error) {
	var total int64
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !(strings.HasSuffix(name, ".m3u8") || strings.HasSuffix(name, ".ts")) {
			continue
		}
		src := filepath.Join(srcDir, name)
		remote := prefix + "/" + name
		upCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		size, err := p.sb.UploadFile(upCtx, p.hlsBucket, remote, src)
		cancel()
		if err != nil {
			return total, fmt.Errorf("upload %s: %w", name, err)
		}
		total += size
	}
	return total, nil
}
