package repository

import (
	"context"

	"github.com/oroya/backend/internal/supabase"
)

type SubscriptionRepository interface {
	Toggle(ctx context.Context, subscriberID, channelID string) (subscribed bool, err error)
	IsSubscribed(ctx context.Context, subscriberID, channelID string) (bool, error)
	CountByChannel(ctx context.Context, channelID string) (int64, error)
}

type subRepo struct{ db *DB }

func NewSubscriptionRepository(db *DB) SubscriptionRepository { return &subRepo{db: db} }

func (r *subRepo) Toggle(ctx context.Context, subscriberID, channelID string) (bool, error) {
	f := supabase.NewFilters()
	f.Set("subscriber_id", "eq."+subscriberID)
	f.Set("channel_id", "eq."+channelID)
	deleted, err := r.db.SB.Delete(ctx, "subscriptions", f)
	if err != nil {
		return false, err
	}
	if deleted > 0 {
		return false, nil
	}
	row := map[string]any{"subscriber_id": subscriberID, "channel_id": channelID}
	if err := r.db.SB.Insert(ctx, "subscriptions", []any{row}, nil); err != nil {
		return false, err
	}
	return true, nil
}

func (r *subRepo) IsSubscribed(ctx context.Context, subscriberID, channelID string) (bool, error) {
	f := supabase.NewFilters()
	f.Set("subscriber_id", "eq."+subscriberID)
	f.Set("channel_id", "eq."+channelID)
	return r.db.SB.Exists(ctx, "subscriptions", f)
}

func (r *subRepo) CountByChannel(ctx context.Context, channelID string) (int64, error) {
	f := supabase.NewFilters()
	f.Set("channel_id", "eq."+channelID)
	return r.db.SB.Count(ctx, "subscriptions", f)
}
