package repository

import (
	"context"
	"strconv"

	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/supabase"
)

type CommentRepository interface {
	Create(ctx context.Context, c *models.Comment) error
	Delete(ctx context.Context, id, userID string) error
	ListByVideo(ctx context.Context, videoID string, limit, offset int) ([]models.Comment, error)
	ToggleLike(ctx context.Context, commentID, userID string) (liked bool, err error)
}

type commentRepo struct{ db *DB }

func NewCommentRepository(db *DB) CommentRepository { return &commentRepo{db: db} }

func (r *commentRepo) Create(ctx context.Context, c *models.Comment) error {
	row := map[string]any{
		"id":       c.ID,
		"video_id": c.VideoID,
		"user_id":  c.UserID,
		"content":  c.Content,
	}
	if c.ParentID != nil && *c.ParentID != "" {
		row["parent_id"] = *c.ParentID
	}
	if !c.CreatedAt.IsZero() {
		row["created_at"] = c.CreatedAt
	}
	return r.db.SB.Insert(ctx, "comments", []any{row}, nil)
}

func (r *commentRepo) Delete(ctx context.Context, id, userID string) error {
	f := supabase.NewFilters()
	f.Set("id", "eq."+id)
	f.Set("user_id", "eq."+userID)
	n, err := r.db.SB.Delete(ctx, "comments", f)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *commentRepo) ListByVideo(ctx context.Context, videoID string, limit, offset int) ([]models.Comment, error) {
	if limit <= 0 {
		limit = 30
	}
	f := supabase.NewFilters()
	f.Set("video_id", "eq."+videoID)
	f.Set("order", "created_at.desc")
	f.Set("limit", strconv.Itoa(limit))
	if offset > 0 {
		f.Set("offset", strconv.Itoa(offset))
	}
	var out []models.Comment
	if err := r.db.SB.Select(ctx, "comments", f, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ToggleLike: trigger keeps comments.likes_count in sync.
func (r *commentRepo) ToggleLike(ctx context.Context, commentID, userID string) (bool, error) {
	f := supabase.NewFilters()
	f.Set("comment_id", "eq."+commentID)
	f.Set("user_id", "eq."+userID)
	deleted, err := r.db.SB.Delete(ctx, "comment_likes", f)
	if err != nil {
		return false, err
	}
	if deleted > 0 {
		return false, nil
	}
	row := map[string]any{"comment_id": commentID, "user_id": userID}
	if err := r.db.SB.Insert(ctx, "comment_likes", []any{row}, nil); err != nil {
		return false, err
	}
	return true, nil
}
