package services

import (
	"context"

	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/repository"
)

type UserService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Get(ctx context.Context, id string) (*models.Profile, error) {
	return s.users.GetProfileByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id string, patch *models.UpdateProfileRequest) (*models.Profile, error) {
	return s.users.UpdateProfile(ctx, id, patch)
}
