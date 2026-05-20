package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/handler"
	"github.com/redis/go-redis/v9"
)

// localLimiter implements a simple token-bucket rate limiter for a single key.
type localLimiter struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

func newLocalLimiter(maxTokens float64, refillRate float64) *localLimiter {
	return &localLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (l *localLimiter) allow() bool {
	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()
	l.tokens += elapsed * l.refillRate
	if l.tokens > l.maxTokens {
		l.tokens = l.maxTokens
	}
	l.lastRefill = now

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// RateLimiter provides Redis-based sliding window rate limiting with an
// in-memory fallback when Redis is unavailable.
type RateLimiter struct {
	client       *redis.Client
	ipLimit      int
	userLimit    int
	windowPeriod time.Duration

	// In-memory fallback limiters (used when Redis is down)
	mu            sync.Mutex
	localLimiters map[string]*localLimiter
}

// NewRateLimiter creates a new rate limiter with the given per-IP and per-user limits.
func NewRateLimiter(client *redis.Client, ipLimit, userLimit int) *RateLimiter {
	return &RateLimiter{
		client:        client,
		ipLimit:       ipLimit,
		userLimit:     userLimit,
		windowPeriod:  time.Minute,
		localLimiters: make(map[string]*localLimiter),
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
		// Redis unavailable — fall back to local in-memory rate limiter
		return rl.allowLocal(key, limit)
	}

	count := cmds[2].(*redis.IntCmd).Val()
	return count <= int64(limit)
}

// allowLocal uses an in-memory token-bucket limiter as a fallback when Redis is down.
// It applies a stricter limit (1/5 of the configured limit) to be conservative.
func (rl *RateLimiter) allowLocal(key string, limit int) bool {
	rl.mu.Lock()
	limiter, exists := rl.localLimiters[key]
	if !exists {
		// Conservative burst: 1/5 of the normal limit
		burst := float64(limit) / 5
		if burst < 1 {
			burst = 1
		}
		// Refill rate: limit tokens per minute → limit/60 tokens per second
		refillRate := float64(limit) / 60.0
		limiter = newLocalLimiter(burst, refillRate)
		rl.localLimiters[key] = limiter
	}
	allowed := limiter.allow()
	rl.mu.Unlock()

	return allowed
}

// MaxBodySize returns middleware that limits the request body to maxBytes.
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
