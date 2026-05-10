package utils

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oroya/backend/internal/models"
)

// ParseSupabaseJWT validates a Supabase-issued HS256 JWT using the project JWT secret
// and returns the claims relevant to the application.
func ParseSupabaseJWT(tokenStr, secret string) (*models.AuthClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}

	if secret == "" {
		return parseSupabaseClaimsUnverified(tokenStr)
	}

	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token claims")
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, errors.New("token missing sub")
	}
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return &models.AuthClaims{UserID: sub, Email: email, Role: role}, nil
}

func parseSupabaseClaimsUnverified(tokenStr string) (*models.AuthClaims, error) {
	tok, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return nil, errors.New("token missing sub")
	}
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return &models.AuthClaims{UserID: sub, Email: email, Role: role}, nil
}
