package repository

import (
	"context"
	"time"

	"github.com/oroya/backend/internal/supabase"
)

type ViewRepository interface {
	// RecordView inserts a row into `views` if no equivalent view from this
	// user (or, for guests, this ip_hash) exists in the past hour. Returns
	// true if a row was inserted (counted). The DB trigger increments
	// videos.views_count on each insert, so the backend does nothing more.
	RecordView(ctx context.Context, videoID string, userID *string, ipHash string) (counted bool, err error)
}

type viewRepo struct{ db *DB }

func NewViewRepository(db *DB) ViewRepository { return &viewRepo{db: db} }

func (r *viewRepo) RecordView(ctx context.Context, videoID string, userID *string, ipHash string) (bool, error) {
	cutoff := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)

	f := supabase.NewFilters()
	f.Set("video_id", "eq."+videoID)
	f.Set("created_at", "gte."+cutoff)
	if userID != nil && *userID != "" {
		f.Set("user_id", "eq."+*userID)
	} else if ipHash != "" {
		f.Set("ip_hash", "eq."+ipHash)
		f.Set("user_id", "is.null")
	} else {
		// Cannot dedup an anonymous view without an IP hash; skip recording
		// so we don't inflate counts from probes/bots.
		return false, nil
	}

	exists, err := r.db.SB.Exists(ctx, "views", f)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	row := map[string]any{"video_id": videoID}
	if userID != nil && *userID != "" {
		row["user_id"] = *userID
	}
	if ipHash != "" {
		row["ip_hash"] = ipHash
	}
	if err := r.db.SB.Insert(ctx, "views", []any{row}, nil); err != nil {
		return false, err
	}
	return true, nil
}
