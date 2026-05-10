package repository

import (
	"context"
	"strconv"
	"strings"

	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/supabase"
)

type SearchRepository interface {
	Videos(ctx context.Context, query string, limit int) ([]models.Video, error)
	Channels(ctx context.Context, query string, limit int) ([]models.Profile, error)
}

type searchRepo struct{ db *DB }

func NewSearchRepository(db *DB) SearchRepository { return &searchRepo{db: db} }

// Videos uses PostgREST's `or=` operator to ILIKE on title/description.
// The `*` wildcard maps to SQL %; commas/parens inside the value must be escaped.
func (r *searchRepo) Videos(ctx context.Context, query string, limit int) ([]models.Video, error) {
	if limit <= 0 {
		limit = 25
	}
	pattern := "*" + sanitizePGRSTLike(query) + "*"

	f := supabase.NewFilters()
	f.Set("status", "eq.ready")
	f.Set("visibility", "eq.public")
	f.Set("or", "(title.ilike."+pattern+",description.ilike."+pattern+")")
	f.Set("order", "views_count.desc,created_at.desc")
	f.Set("limit", strconv.Itoa(limit))

	var out []models.Video
	if err := r.db.SB.Select(ctx, "videos", f, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *searchRepo) Channels(ctx context.Context, query string, limit int) ([]models.Profile, error) {
	if limit <= 0 {
		limit = 25
	}
	pattern := "*" + sanitizePGRSTLike(query) + "*"

	f := supabase.NewFilters()
	f.Set("or", "(username.ilike."+pattern+",display_name.ilike."+pattern+")")
	f.Set("limit", strconv.Itoa(limit))

	var out []models.Profile
	if err := r.db.SB.Select(ctx, "profiles", f, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PostgREST treats comma, parentheses, and dot as syntactic. Strip them out
// of user-provided search terms (callers should also length-cap the query).
func sanitizePGRSTLike(s string) string {
	r := strings.NewReplacer(",", " ", "(", " ", ")", " ", "*", " ", ".", " ")
	return strings.TrimSpace(r.Replace(s))
}
