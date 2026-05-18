package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/handler"
	"github.com/redis/go-redis/v9"
)

// RateLimiter provides Redis-based sliding window rate limiting.
type RateLimiter struct {
	client       *redis.Client
	ipLimit      int
	userLimit    int
	windowPeriod time.Duration
}

// NewRateLimiter creates a new rate limiter with the given per-IP and per-user limits.
func NewRateLimiter(client *redis.Client, ipLimit, userLimit int) *RateLimiter {
	return &RateLimiter{
		client:       client,
		ipLimit:      ipLimit,
		userLimit:    userLimit,
		windowPeriod: time.Minute,
	}
}

// IPMiddleware returns an HTTP middleware that enforces IP-based rate limiting only.
// This should be applied globally before auth middleware.
func (rl *RateLimiter) IPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Rate limit by IP
		ip := r.RemoteAddr
		key := fmt.Sprintf("ratelimit:ip:%s", ip)
		if !rl.allow(ctx, key, rl.ipLimit) {
			handler.RespondError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// UserMiddleware returns an HTTP middleware that enforces per-user rate limiting.
// This should be applied after JWT auth middleware so that user claims are available.
func (rl *RateLimiter) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if claims, ok := GetUserClaims(ctx); ok {
			userKey := fmt.Sprintf("ratelimit:user:%s", claims.UserID.String())
			if !rl.allow(ctx, userKey, rl.userLimit) {
				handler.RespondError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ctx context.Context, key string, limit int) bool {
	now := time.Now()
	windowStart := now.Add(-rl.windowPeriod)

	pipe := rl.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now.UnixNano()), Member: now.UnixNano()})
	pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, rl.windowPeriod)

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		// If Redis is unavailable, allow the request
		return true
	}

	count := cmds[2].(*redis.IntCmd).Val()
	return count <= int64(limit)
}
