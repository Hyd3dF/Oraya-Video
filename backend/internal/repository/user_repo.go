package repository

import (
	"context"

	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/supabase"
)

type UserRepository interface {
	CreateProfile(ctx context.Context, p *models.Profile) error
	GetProfileByID(ctx context.Context, id string) (*models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (*models.Profile, error)
	UpdateProfile(ctx context.Context, id string, patch *models.UpdateProfileRequest) (*models.Profile, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
}

type userRepo struct{ db *DB }

func NewUserRepository(db *DB) UserRepository { return &userRepo{db: db} }

func (r *userRepo) CreateProfile(ctx context.Context, p *models.Profile) error {
	row := map[string]any{
		"id":           p.ID,
		"username":     p.Username,
		"email":        p.Email,
		"login_type":   p.LoginType,
	}
	setIfNotEmpty(row, "real_name", p.RealName)
	setIfNotEmpty(row, "display_name", p.DisplayName)
	setIfNotEmpty(row, "avatar_url", p.AvatarURL)
	setIfNotEmpty(row, "banner_url", p.BannerURL)
	setIfNotEmpty(row, "bio", p.Bio)
	if !p.CreatedAt.IsZero() {
		row["created_at"] = p.CreatedAt
	}
	if !p.UpdatedAt.IsZero() {
		row["updated_at"] = p.UpdatedAt
	}
	return r.db.SB.Insert(ctx, "profiles", []any{row}, nil)
}

func (r *userRepo) GetProfileByID(ctx context.Context, id string) (*models.Profile, error) {
	var p models.Profile
	if err := r.db.SB.SelectOne(ctx, "profiles", idFilter(id), &p); err != nil {
		return nil, translateNotFound(err)
	}
	return &p, nil
}

func (r *userRepo) GetProfileByUsername(ctx context.Context, username string) (*models.Profile, error) {
	f := supabase.NewFilters()
	f.Set("username", "eq."+username)
	var p models.Profile
	if err := r.db.SB.SelectOne(ctx, "profiles", f, &p); err != nil {
		return nil, translateNotFound(err)
	}
	return &p, nil
}

func (r *userRepo) UpdateProfile(ctx context.Context, id string, patch *models.UpdateProfileRequest) (*models.Profile, error) {
	body := map[string]any{"updated_at": "now()"}
	setIfPtr(body, "real_name", patch.RealName)
	setIfPtr(body, "username", patch.Username)
	setIfPtr(body, "display_name", patch.DisplayName)
	setIfPtr(body, "avatar_url", patch.AvatarURL)
	setIfPtr(body, "banner_url", patch.BannerURL)
	setIfPtr(body, "bio", patch.Bio)

	var rows []models.Profile
	if err := r.db.SB.Update(ctx, "profiles", idFilter(id), body, &rows); err != nil {
		return nil, translateNotFound(err)
	}
	if len(rows) == 0 {
		return nil, ErrNotFound
	}
	return &rows[0], nil
}

func (r *userRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	f := supabase.NewFilters()
	f.Set("username", "eq."+username)
	return r.db.SB.Exists(ctx, "profiles", f)
}

func setIfNotEmpty(m map[string]any, key, val string) {
	if val != "" {
		m[key] = val
	}
}

func setIfPtr(m map[string]any, key string, val *string) {
	if val != nil {
		m[key] = *val
	}
}
