package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/repository"
	"github.com/oroya/backend/internal/supabase"
	"github.com/oroya/backend/internal/worker"
)

var ErrForbidden = errors.New("forbidden")

type VideoService struct {
	videos        repository.VideoRepository
	views         repository.ViewRepository
	queue         worker.Queue
	sb            *supabase.Client
	sourceBucket  string
	processVideos bool
}

func NewVideoService(videos repository.VideoRepository, views repository.ViewRepository, queue worker.Queue, sb *supabase.Client, sourceBucket string, processVideos bool) *VideoService {
	return &VideoService{
		videos:        videos,
		views:         views,
		queue:         queue,
		sb:            sb,
		sourceBucket:  sourceBucket,
		processVideos: processVideos,
	}
}

func (s *VideoService) Create(ctx context.Context, ownerID string, req *models.CreateVideoRequest) (*models.Video, error) {
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.StoragePath) == "" {
		return nil, ErrInvalidInput
	}
	visibility := req.Visibility
	if visibility == "" {
		visibility = "public"
	}

	now := time.Now().UTC()
	v := &models.Video{
		ID:              uuid.NewString(),
		OwnerID:         ownerID,
		Title:           req.Title,
		Description:     req.Description,
		ThumbnailURL:    req.ThumbnailURL,
		DurationSeconds: req.DurationSeconds,
		Visibility:      visibility,
		Status:          "ready",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if s.processVideos {
		v.Status = "processing"
	}
	if err := s.videos.Create(ctx, v); err != nil {
		return nil, err
	}

	if !s.processVideos {
		if s.sb == nil || strings.TrimSpace(s.sourceBucket) == "" {
			return nil, errors.New("source video storage is not configured")
		}
		asset := &models.VideoAsset{
			ID:          uuid.NewString(),
			VideoID:     v.ID,
			Quality:     "source",
			PlaylistURL: s.sb.PublicURL(s.sourceBucket, req.StoragePath),
			CreatedAt:   now,
		}
		if err := s.videos.AddAsset(ctx, asset); err != nil {
			return nil, err
		}
		return v, nil
	}

	if err := s.queue.Enqueue(worker.Job{
		VideoID:     v.ID,
		StoragePath: req.StoragePath,
		OwnerID:     ownerID,
	}); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *VideoService) Get(ctx context.Context, id string) (*models.Video, []models.VideoAsset, []models.VideoLink, error) {
	v, err := s.videos.GetByID(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	assets, err := s.videos.ListAssets(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	links, err := s.videos.ListLinks(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	return v, assets, links, nil
}

func (s *VideoService) List(ctx context.Context, limit, offset int) ([]models.Video, error) {
	if limit <= 0 || limit > 100 {
		limit = 24
	}
	return s.videos.List(ctx, repository.ListVideosOpts{Limit: limit, Offset: offset, Status: "ready"})
}

func (s *VideoService) ListByOwner(ctx context.Context, ownerID string, limit, offset int) ([]models.Video, error) {
	if strings.TrimSpace(ownerID) == "" {
		return nil, ErrInvalidInput
	}
	if limit <= 0 || limit > 100 {
		limit = 24
	}
	return s.videos.List(ctx, repository.ListVideosOpts{Limit: limit, Offset: offset, OwnerID: ownerID})
}

func (s *VideoService) Update(ctx context.Context, id, userID string, patch *models.UpdateVideoRequest) (*models.Video, error) {
	v, err := s.videos.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if v.OwnerID != userID {
		return nil, ErrForbidden
	}
	return s.videos.Update(ctx, id, patch)
}

func (s *VideoService) Delete(ctx context.Context, id, userID string) error {
	v, err := s.videos.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if v.OwnerID != userID {
		return ErrForbidden
	}
	return s.videos.Delete(ctx, id)
}

// RegisterView records a view, deduplicating by user/IP within the past hour.
// The DB trigger increments videos.views_count for each inserted row.
func (s *VideoService) RegisterView(ctx context.Context, videoID string, userID *string, ipHash string) error {
	_, err := s.views.RecordView(ctx, videoID, userID, ipHash)
	return err
}

func (s *VideoService) ToggleLike(ctx context.Context, videoID, userID string) (bool, error) {
	return s.videos.ToggleLike(ctx, videoID, userID)
}

func (s *VideoService) AddLink(ctx context.Context, videoID, userID string, req *models.CreateVideoLinkRequest) (*models.VideoLink, error) {
	v, err := s.videos.GetByID(ctx, videoID)
	if err != nil {
		return nil, err
	}
	if v.OwnerID != userID {
		return nil, ErrForbidden
	}
	link := &models.VideoLink{
		ID:        uuid.NewString(),
		VideoID:   videoID,
		Title:     req.Title,
		URL:       req.URL,
		SortOrder: req.SortOrder,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.videos.AddLink(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *VideoService) DeleteLink(ctx context.Context, videoID, linkID, userID string) error {
	v, err := s.videos.GetByID(ctx, videoID)
	if err != nil {
		return err
	}
	if v.OwnerID != userID {
		return ErrForbidden
	}
	return s.videos.DeleteLink(ctx, videoID, linkID)
}

func (s *VideoService) ListLinks(ctx context.Context, videoID string) ([]models.VideoLink, error) {
	return s.videos.ListLinks(ctx, videoID)
}
