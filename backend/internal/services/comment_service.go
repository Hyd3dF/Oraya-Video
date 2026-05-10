package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/repository"
)

type CommentService struct {
	comments repository.CommentRepository
}

func NewCommentService(comments repository.CommentRepository) *CommentService {
	return &CommentService{comments: comments}
}

func (s *CommentService) Create(ctx context.Context, videoID, userID string, req *models.CreateCommentRequest) (*models.Comment, error) {
	if strings.TrimSpace(req.Content) == "" {
		return nil, ErrInvalidInput
	}
	c := &models.Comment{
		ID:        uuid.NewString(),
		VideoID:   videoID,
		UserID:    userID,
		ParentID:  req.ParentID,
		Content:   req.Content,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.comments.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CommentService) Delete(ctx context.Context, id, userID string) error {
	return s.comments.Delete(ctx, id, userID)
}

func (s *CommentService) List(ctx context.Context, videoID string, limit, offset int) ([]models.Comment, error) {
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	return s.comments.ListByVideo(ctx, videoID, limit, offset)
}

func (s *CommentService) ToggleLike(ctx context.Context, commentID, userID string) (bool, error) {
	return s.comments.ToggleLike(ctx, commentID, userID)
}
