package services

import (
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/oroya/backend/internal/models"
	"github.com/oroya/backend/internal/supabase"
)

type UploadService struct {
	sb        *supabase.Client
	rawBucket string
}

func NewUploadService(sb *supabase.Client, rawBucket string) *UploadService {
	return &UploadService{sb: sb, rawBucket: rawBucket}
}

// SignUploadURL generates a signed Supabase Storage URL where the frontend can
// PUT the raw video file directly. The backend never touches the bytes.
func (s *UploadService) SignUploadURL(userID, filename string) (*models.UploadURLResponse, error) {
	id := uuid.NewString()
	ext := path.Ext(filename)
	if ext == "" {
		ext = ".mp4"
	}
	storagePath := path.Join("users", userID, id+ext)

	signed, err := s.sb.CreateSignedUploadURL(s.rawBucket, storagePath)
	if err != nil {
		return nil, err
	}
	return &models.UploadURLResponse{
		UploadURL:   signed.URL,
		StoragePath: storagePath,
		Token:       signed.Token,
		ExpiresAt:   time.Now().Add(2 * time.Hour).Unix(),
	}, nil
}
