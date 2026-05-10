package services

import (
	"context"
	"strings"

	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/repository"
)

type SearchResults struct {
	Query    string           `json:"query"`
	Videos   []models.Video   `json:"videos"`
	Channels []models.Profile `json:"channels"`
}

type SearchService struct {
	repo repository.SearchRepository
}

func NewSearchService(repo repository.SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) Search(ctx context.Context, query string, limit int) (*SearchResults, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, ErrInvalidInput
	}
	vids, err := s.repo.Videos(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	chans, err := s.repo.Channels(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	return &SearchResults{Query: q, Videos: vids, Channels: chans}, nil
}
