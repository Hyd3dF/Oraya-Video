package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/repository"
	"github.com/oroya/backend/internal/supabase"
	"github.com/oroya/backend/internal/utils"
)

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrUsernameTaken = errors.New("username already taken")
)

type AuthService struct {
	users repository.UserRepository
	sb    *supabase.Client
}

func NewAuthService(users repository.UserRepository, sb *supabase.Client) *AuthService {
	return &AuthService{users: users, sb: sb}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if !utils.ValidEmail(req.Email) || !utils.ValidPassword(req.Password) || !utils.ValidUsername(req.Username) || strings.TrimSpace(req.RealName) == "" {
		return nil, ErrInvalidInput
	}

	taken, err := s.users.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrUsernameTaken
	}

	authUser, err := s.sb.SignUp(req.Email, req.Password, map[string]any{
		"username":  req.Username,
		"real_name": req.RealName,
	})
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	profile := &models.Profile{
		ID:        authUser.ID,
		RealName:  req.RealName,
		Username:  req.Username,
		Email:     req.Email,
		LoginType: "email",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.users.CreateProfile(ctx, profile); err != nil {
		return nil, err
	}

	session, err := s.sb.PasswordSignIn(req.Email, req.Password)
	if err != nil {
		return &models.AuthResponse{
			TokenType: "pending_confirmation",
			User:      *profile,
		}, nil
	}
	return sessionToResponse(session, profile), nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !utils.ValidEmail(req.Email) || req.Password == "" {
		return nil, ErrInvalidInput
	}
	session, err := s.sb.PasswordSignIn(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	profile, err := s.users.GetProfileByID(ctx, session.User.ID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	return sessionToResponse(session, profile), nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, ErrInvalidInput
	}
	session, err := s.sb.RefreshSession(refreshToken)
	if err != nil {
		return nil, err
	}
	profile, _ := s.users.GetProfileByID(ctx, session.User.ID)
	return sessionToResponse(session, profile), nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	return s.sb.SignOut(accessToken)
}

func (s *AuthService) Me(ctx context.Context, userID string) (*models.Profile, error) {
	return s.users.GetProfileByID(ctx, userID)
}

func sessionToResponse(s *supabase.GoTrueSession, profile *models.Profile) *models.AuthResponse {
	resp := &models.AuthResponse{
		AccessToken:  s.AccessToken,
		RefreshToken: s.RefreshToken,
		ExpiresIn:    s.ExpiresIn,
		TokenType:    s.TokenType,
	}
	if profile != nil {
		resp.User = *profile
	} else {
		resp.User = models.Profile{ID: s.User.ID, Email: s.User.Email}
	}
	return resp
}
