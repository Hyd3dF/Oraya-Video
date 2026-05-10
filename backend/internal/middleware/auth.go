package middleware

import (
	"net/http"
	"strings"

	"github.com/oroya/backend/internal/utils"
)

// RequireAuth validates a Supabase JWT and injects claims into the request context.
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok := bearerToken(r)
			if tok == "" {
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing bearer token")
				return
			}
			claims, err := utils.ParseSupabaseJWT(tok, jwtSecret)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid token")
				return
			}
			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth attaches claims if a valid token is present, otherwise passes through.
func OptionalAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tok := bearerToken(r); tok != "" {
				if claims, err := utils.ParseSupabaseJWT(tok, jwtSecret); err == nil {
					r = r.WithContext(WithClaims(r.Context(), claims))
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if h == "" {
		return ""
	}
	const p = "Bearer "
	if !strings.HasPrefix(h, p) {
		return ""
	}
	return strings.TrimSpace(h[len(p):])
}
