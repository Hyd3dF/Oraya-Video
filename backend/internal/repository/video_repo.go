package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/supabase"
)

type ListVideosOpts struct {
	Limit   int
	Offset  int
	Status  string
	OwnerID string
}

type VideoRepository interface {
	Create(ctx context.Context, v *models.Video) error
	GetByID(ctx context.Context, id string) (*models.Video, error)
	List(ctx context.Context, opts ListVideosOpts) ([]models.Video, error)
	Update(ctx context.Context, id string, patch *models.UpdateVideoRequest) (*models.Video, error)
	Delete(ctx context.Context, id string) error
	SetStatus(ctx context.Context, id, status string) error

	AddAsset(ctx context.Context, a *models.VideoAsset) error
	ListAssets(ctx context.Context, videoID string) ([]models.VideoAsset, error)

	AddLink(ctx context.Context, l *models.VideoLink) error
	DeleteLink(ctx context.Context, videoID, linkID string) error
	ListLinks(ctx context.Context, videoID string) ([]models.VideoLink, error)

	// ToggleLike inserts or removes a row in video_likes; the trigger keeps
	// videos.likes_count in sync. Returns true if the user now likes the video.
	ToggleLike(ctx context.Context, videoID, userID string) (liked bool, err error)
}

type videoRepo struct{ db *DB }

func NewVideoRepository(db *DB) VideoRepository { return &videoRepo{db: db} }

func (r *videoRepo) Create(ctx context.Context, v *models.Video) error {
	row := map[string]any{
		"id":         v.ID,
		"owner_id":   v.OwnerID,
		"title":      v.Title,
		"visibility": v.Visibility,
		"status":     v.Status,
	}
	setIfNotEmpty(row, "description", v.Description)
	setIfNotEmpty(row, "thumbnail_url", v.ThumbnailURL)
	if v.DurationSeconds > 0 {
		row["duration_seconds"] = v.DurationSeconds
	}
	if !v.CreatedAt.IsZero() {
		row["created_at"] = v.CreatedAt
	}
	if !v.UpdatedAt.IsZero() {
		row["updated_at"] = v.UpdatedAt
	}
	return r.db.SB.Insert(ctx, "videos", []any{row}, nil)
}

func (r *videoRepo) GetByID(ctx context.Context, id string) (*models.Video, error) {
	var v models.Video
	if err := r.db.SB.SelectOne(ctx, "videos", idFilter(id), &v); err != nil {
		return nil, translateNotFound(err)
	}
	return &v, nil
}

func (r *videoRepo) List(ctx context.Context, opts ListVideosOpts) ([]models.Video, error) {
	if opts.Limit <= 0 {
		opts.Limit = 24
	}
	f := supabase.NewFilters()
	if opts.Status != "" {
		f.Set("status", "eq."+opts.Status)
	}
	if opts.OwnerID != "" {
		f.Set("owner_id", "eq."+opts.OwnerID)
	}
	f.Set("order", "created_at.desc")
	f.Set("limit", strconv.Itoa(opts.Limit))
	if opts.Offset > 0 {
		f.Set("offset", strconv.Itoa(opts.Offset))
	}
	var out []models.Video
	if err := r.db.SB.Select(ctx, "videos", f, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *videoRepo) Update(ctx context.Context, id string, patch *models.UpdateVideoRequest) (*models.Video, error) {
	body := map[string]any{"updated_at": time.Now().UTC()}
	setIfPtr(body, "title", patch.Title)
	setIfPtr(body, "description", patch.Description)
	setIfPtr(body, "visibility", patch.Visibility)

	var rows []models.Video
	if err := r.db.SB.Update(ctx, "videos", idFilter(id), body, &rows); err != nil {
		return nil, translateNotFound(err)
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	return &rows[0], nil
}

func (r *videoRepo) Delete(ctx context.Context, id string) error {
	n, err := r.db.SB.Delete(ctx, "videos", idFilter(id))
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *videoRepo) SetStatus(ctx context.Context, id, status string) error {
	body := map[string]any{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}
	return r.db.SB.Update(ctx, "videos", idFilter(id), body, nil)
}

// --- Assets ---

func (r *videoRepo) AddAsset(ctx context.Context, a *models.VideoAsset) error {
	row := map[string]any{
		"id":           a.ID,
		"video_id":     a.VideoID,
		"quality":      a.Quality,
		"playlist_url": a.PlaylistURL,
	}
	setIfNotEmpty(row, "master_url", a.MasterURL)
	if a.Width > 0 {
		row["width"] = a.Width
	}
	if a.Height > 0 {
		row["height"] = a.Height
	}
	if a.Bitrate > 0 {
		row["bitrate"] = a.Bitrate
	}
	if a.SizeBytes > 0 {
		row["size_bytes"] = a.SizeBytes
	}
	if !a.CreatedAt.IsZero() {
		row["created_at"] = a.CreatedAt
	}
	return r.db.SB.Upsert(ctx, "video_assets", []any{row}, "video_id,quality", nil)
}

func (r *videoRepo) ListAssets(ctx context.Context, videoID string) ([]models.VideoAsset, error) {
	f := supabase.NewFilters()
	f.Set("video_id", "eq."+videoID)
	f.Set("order", "bitrate.desc.nullslast")
	var out []models.VideoAsset
	if err := r.db.SB.Select(ctx, "video_assets", f, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// --- Links ---

func (r *videoRepo) AddLink(ctx context.Context, l *models.VideoLink) error {
	row := map[string]any{
		"id":         l.ID,
		"video_id":   l.VideoID,
		"url":        l.URL,
		"sort_order": l.SortOrder,
	}
	setIfNotEmpty(row, "title", l.Title)
	if !l.CreatedAt.IsZero() {
		row["created_at"] = l.CreatedAt
	}
	return r.db.SB.Insert(ctx, "video_links", []any{row}, nil)
}

func (r *videoRepo) DeleteLink(ctx context.Context, videoID, linkID string) error {
	f := supabase.NewFilters()
	f.Set("id", "eq."+linkID)
	f.Set("video_id", "eq."+videoID)
	n, err := r.db.SB.Delete(ctx, "video_links", f)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *videoRepo) ListLinks(ctx context.Context, videoID string) ([]models.VideoLink, error) {
	f := supabase.NewFilters()
	f.Set("video_id", "eq."+videoID)
	f.Set("order", "sort_order.asc,created_at.asc")
	var out []models.VideoLink
	if err := r.db.SB.Select(ctx, "video_links", f, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// --- Likes ---
//
// Toggle is best-effort: try DELETE first; if nothing was deleted, INSERT.
// The likes_count column is updated by the database trigger in 010_triggers.sql.
func (r *videoRepo) ToggleLike(ctx context.Context, videoID, userID string) (bool, error) {
	f := supabase.NewFilters()
	f.Set("video_id", "eq."+videoID)
	f.Set("user_id", "eq."+userID)

	deleted, err := r.db.SB.Delete(ctx, "video_likes", f)
	if err != nil {
		return false, err
	}
	if deleted > 0 {
		return false, nil
	}
	row := map[string]any{"video_id": videoID, "user_id": userID}
	if err := r.db.SB.Insert(ctx, "video_likes", []any{row}, nil); err != nil {
		return false, err
	}
	return true, nil
}
