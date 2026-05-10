package services

import (
	"context"

	"github.com/oroya/backend/internal/repository"
	"github.com/oroya/backend/internal/supabase"
	"github.com/oroya/backend/internal/worker"
)

type AdminService struct {
	repo      repository.AdminRepository
	processor *worker.Processor
	sb        *supabase.Client
	buckets   AdminBuckets
	db        *repository.DB
}

type AdminBuckets struct {
	Raw        string
	HLS        string
	Thumbnails string
}

func NewAdminService(repo repository.AdminRepository, processor *worker.Processor, sb *supabase.Client, buckets AdminBuckets, db *repository.DB) *AdminService {
	return &AdminService{repo: repo, processor: processor, sb: sb, buckets: buckets, db: db}
}

func (s *AdminService) Stats(ctx context.Context) (*repository.AdminStats, error) {
	return s.repo.Stats(ctx)
}

type QueueStatus struct {
	Pending    int            `json:"pending"`
	Processing int64          `json:"processing"`
	Completed  int64          `json:"completed"`
	Failed     int64          `json:"failed"`
	Statuses   map[string]int64 `json:"video_statuses"`
}

func (s *AdminService) Queue(ctx context.Context) (*QueueStatus, error) {
	statuses, err := s.repo.VideoStatusCounts(ctx)
	if err != nil {
		return nil, err
	}
	q := &QueueStatus{Statuses: statuses}
	if s.processor != nil {
		st := s.processor.Status()
		q.Pending = st.Pending
		q.Processing = st.Processing
		q.Completed = st.Completed
		q.Failed = st.Failed
	}
	return q, nil
}

func (s *AdminService) WorkerStatus() *worker.Status {
	if s.processor == nil {
		return nil
	}
	st := s.processor.Status()
	return &st
}

type StorageStatus struct {
	Buckets map[string]*supabase.BucketUsage `json:"buckets"`
}

func (s *AdminService) StorageStatus(ctx context.Context) (*StorageStatus, error) {
	out := &StorageStatus{Buckets: map[string]*supabase.BucketUsage{}}
	for label, name := range map[string]string{
		"raw":        s.buckets.Raw,
		"hls":        s.buckets.HLS,
		"thumbnails": s.buckets.Thumbnails,
	} {
		usage, err := s.sb.BucketUsage(ctx, name)
		if err != nil {
			out.Buckets[label] = &supabase.BucketUsage{Bucket: name}
			continue
		}
		out.Buckets[label] = usage
	}
	return out, nil
}

type Health struct {
	Status   string `json:"status"`
	Database bool   `json:"database"`
	Storage  bool   `json:"storage"`
	Worker   bool   `json:"worker"`
}

func (s *AdminService) Health(ctx context.Context) *Health {
	h := &Health{Status: "ok"}
	if s.db != nil && s.db.Ping(ctx) == nil {
		h.Database = true
	}
	if _, err := s.sb.BucketUsage(ctx, s.buckets.HLS); err == nil {
		h.Storage = true
	}
	if s.processor != nil {
		h.Worker = true
	}
	if !h.Database || !h.Storage {
		h.Status = "degraded"
	}
	return h
}
