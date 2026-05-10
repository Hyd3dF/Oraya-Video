package models

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RealName string `json:"real_name"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int     `json:"expires_in"`
	TokenType    string  `json:"token_type"`
	User         Profile `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

type CreateVideoRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Visibility      string `json:"visibility"`
	StoragePath     string `json:"storage_path"`
	DurationSeconds int    `json:"duration_seconds"`
	ThumbnailURL    string `json:"thumbnail_url"`
}

type UpdateVideoRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
}

type UploadURLResponse struct {
	UploadURL   string `json:"upload_url"`
	StoragePath string `json:"storage_path"`
	Token       string `json:"token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type CreateCommentRequest struct {
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id,omitempty"`
}

type CreateVideoLinkRequest struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

type UpdateProfileRequest struct {
	RealName    *string `json:"real_name,omitempty"`
	Username    *string `json:"username,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	BannerURL   *string `json:"banner_url,omitempty"`
	Bio         *string `json:"bio,omitempty"`
}
