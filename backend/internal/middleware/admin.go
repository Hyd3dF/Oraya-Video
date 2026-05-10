package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/oroya/backend/internal/utils"
)

// RequireAdmin gates admin endpoints behind a static API token.
// In production this should be replaced or augmented with a role claim check.
func RequireAdmin(adminToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if adminToken == "" {
				utils.WriteError(w, http.StatusServiceUnavailable, "admin_disabled", "admin token not configured")
				return
			}
			supplied := r.Header.Get("X-Admin-Token")
			if subtle.ConstantTimeCompare([]byte(supplied), []byte(adminToken)) != 1 {
				utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid admin token")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
