package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/oroya/backend/internal/utils"
)

// Simple in-memory token bucket per client IP.
// For production, swap with a Redis-backed limiter.
type bucket struct {
	tokens     float64
	lastRefill time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rps     float64
	burst   float64
}

func NewRateLimiter(rps, burst int) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*bucket),
		rps:     float64(rps),
		burst:   float64(burst),
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[key]
	if !ok {
		rl.buckets[key] = &bucket{tokens: rl.burst - 1, lastRefill: now}
		return true
	}
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens = min(rl.burst, b.tokens+elapsed*rl.rps)
	b.lastRefill = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r)
		if !rl.allow(key) {
			utils.WriteError(w, http.StatusTooManyRequests, "rate_limited", "too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

