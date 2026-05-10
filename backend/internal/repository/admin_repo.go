package repository

import (
	"context"

	"github.com/oroya/backend/internal/supabase"
)

type AdminStats struct {
	UsersTotal       int64 `json:"users_total"`
	VideosTotal      int64 `json:"videos_total"`
	VideosProcessing int64 `json:"videos_processing"`
	VideosReady      int64 `json:"videos_ready"`
	VideosFailed     int64 `json:"videos_failed"`
	CommentsTotal    int64 `json:"comments_total"`
	ViewsTotal       int64 `json:"views_total"`
	LikesTotal       int64 `json:"likes_total"`
}

type AdminRepository interface {
	Stats(ctx context.Context) (*AdminStats, error)
	VideoStatusCounts(ctx context.Context) (map[string]int64, error)
}

type adminRepo struct{ db *DB }

func NewAdminRepository(db *DB) AdminRepository { return &adminRepo{db: db} }

func (r *adminRepo) Stats(ctx context.Context) (*AdminStats, error) {
	s := &AdminStats{}
	type job struct {
		dst    *int64
		table  string
		filter string
		val    string
	}
	jobs := []job{
		{&s.UsersTotal, "profiles", "", ""},
		{&s.VideosTotal, "videos", "", ""},
		{&s.VideosProcessing, "videos", "status", "eq.processing"},
		{&s.VideosReady, "videos", "status", "eq.ready"},
		{&s.VideosFailed, "videos", "status", "eq.failed"},
		{&s.CommentsTotal, "comments", "", ""},
		{&s.ViewsTotal, "views", "", ""},
	}
	for _, j := range jobs {
		f := supabase.NewFilters()
		if j.filter != "" {
			f.Set(j.filter, j.val)
		}
		n, err := r.db.SB.Count(ctx, j.table, f)
		if err != nil {
			return nil, err
		}
		*j.dst = n
	}
	vl, err := r.db.SB.Count(ctx, "video_likes", supabase.NewFilters())
	if err != nil {
		return nil, err
	}
	cl, err := r.db.SB.Count(ctx, "comment_likes", supabase.NewFilters())
	if err != nil {
		return nil, err
	}
	s.LikesTotal = vl + cl
	return s, nil
}

func (r *adminRepo) VideoStatusCounts(ctx context.Context) (map[string]int64, error) {
	out := map[string]int64{}
	for _, status := range []string{"processing", "ready", "failed"} {
		f := supabase.NewFilters()
		f.Set("status", "eq."+status)
		n, err := r.db.SB.Count(ctx, "videos", f)
		if err != nil {
			return nil, err
		}
		out[status] = n
	}
	return out, nil
}
