package services

import (
	"context"
	"errors"

	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/repository"
)

var ErrSelfSubscribe = errors.New("cannot subscribe to own channel")

type ChannelView struct {
	Profile         *models.Profile `json:"profile"`
	Videos          []models.Video  `json:"videos"`
	SubscriberCount int64           `json:"subscriber_count"`
	IsSubscribed    bool            `json:"is_subscribed"`
}

type ChannelService struct {
	users repository.UserRepository
	vids  repository.VideoRepository
	subs  repository.SubscriptionRepository
}

func NewChannelService(users repository.UserRepository, vids repository.VideoRepository, subs repository.SubscriptionRepository) *ChannelService {
	return &ChannelService{users: users, vids: vids, subs: subs}
}

// Get returns the public channel page payload. viewerID is the authenticated
// user's id (empty for guests) and controls the IsSubscribed field.
func (s *ChannelService) Get(ctx context.Context, channelID, viewerID string) (*ChannelView, error) {
	profile, err := s.users.GetProfileByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	vids, err := s.vids.List(ctx, repository.ListVideosOpts{
		Limit:   24,
		Status:  "ready",
		OwnerID: channelID,
	})
	if err != nil {
		return nil, err
	}
	count, err := s.subs.CountByChannel(ctx, channelID)
	if err != nil {
		return nil, err
	}
	view := &ChannelView{
		Profile:         profile,
		Videos:          vids,
		SubscriberCount: count,
	}
	if viewerID != "" && viewerID != channelID {
		view.IsSubscribed, _ = s.subs.IsSubscribed(ctx, viewerID, channelID)
	}
	return view, nil
}

func (s *ChannelService) ToggleSubscribe(ctx context.Context, subscriberID, channelID string) (bool, error) {
	if subscriberID == channelID {
		return false, ErrSelfSubscribe
	}
	if _, err := s.users.GetProfileByID(ctx, channelID); err != nil {
		return false, err
	}
	return s.subs.Toggle(ctx, subscriberID, channelID)
}
