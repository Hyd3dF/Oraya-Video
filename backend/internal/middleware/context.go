package middleware

import (
	"context"

	"github.com/oroya/backend/internal/models"
)

type ctxKey string

const claimsKey ctxKey = "auth.claims"

func WithClaims(ctx context.Context, c *models.AuthClaims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

func ClaimsFrom(ctx context.Context) (*models.AuthClaims, bool) {
	c, ok := ctx.Value(claimsKey).(*models.AuthClaims)
	return c, ok && c != nil
}
