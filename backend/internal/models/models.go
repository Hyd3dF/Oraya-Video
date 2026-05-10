package models

import "time"

type Profile struct {
	ID          string    `json:"id"`
	RealName    string    `json:"real_name"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name,omitempty"`
	Email       string    `json:"email"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	BannerURL   string    `json:"banner_url,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	LoginType   string    `json:"login_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Video struct {
	ID              string    `json:"id"`
	OwnerID         string    `json:"owner_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	ThumbnailURL    string    `json:"thumbnail_url,omitempty"`
	DurationSeconds int       `json:"duration_seconds"`
	ViewsCount      int64     `json:"views_count"`
	LikesCount      int64     `json:"likes_count"`
	Visibility      string    `json:"visibility"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type VideoAsset struct {
	ID          string    `json:"id"`
	VideoID     string    `json:"video_id"`
	Quality     string    `json:"quality"`
	PlaylistURL string    `json:"playlist_url"`
	MasterURL   string    `json:"master_url,omitempty"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Bitrate     int       `json:"bitrate"`
	SizeBytes   int64     `json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`
}

type VideoLink struct {
	ID        string    `json:"id"`
	VideoID   string    `json:"video_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID         string    `json:"id"`
	VideoID    string    `json:"video_id"`
	UserID     string    `json:"user_id"`
	ParentID   *string   `json:"parent_id,omitempty"`
	Content    string    `json:"content"`
	LikesCount int64     `json:"likes_count"`
	CreatedAt  time.Time `json:"created_at"`
}

type Subscription struct {
	ID           int64     `json:"id"`
	SubscriberID string    `json:"subscriber_id"`
	ChannelID    string    `json:"channel_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type View struct {
	ID        int64     `json:"id"`
	VideoID   string    `json:"video_id"`
	UserID    *string   `json:"user_id,omitempty"`
	IPHash    string    `json:"ip_hash,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthClaims represents validated Supabase JWT claims attached to request context.
type AuthClaims struct {
	UserID string
	Email  string
	Role   string
}
