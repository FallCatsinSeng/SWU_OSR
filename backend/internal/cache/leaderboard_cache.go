package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// leaderboardTTL is the cache TTL for leaderboard data.
	// Set slightly shorter than refresh interval so stale data is rare.
	leaderboardTTL = 14 * time.Minute

	// userSummaryTTL is the cache TTL for individual user summaries.
	userSummaryTTL = 10 * time.Minute
)

// CachedLeaderboardService wraps a domain.LeaderboardService with Redis caching.
type CachedLeaderboardService struct {
	inner  domain.LeaderboardService
	rdb    *redis.Client
	logger *zap.Logger
}

// NewCachedLeaderboardService creates a caching decorator around a LeaderboardService.
func NewCachedLeaderboardService(inner domain.LeaderboardService, rdb *redis.Client, logger *zap.Logger) domain.LeaderboardService {
	return &CachedLeaderboardService{
		inner:  inner,
		rdb:    rdb,
		logger: logger,
	}
}

// GetLeaderboard attempts to serve from cache, falling back to the inner service.
func (c *CachedLeaderboardService) GetLeaderboard(ctx context.Context, period domain.LeaderboardPeriod, limit, offset int) (*domain.LeaderboardResult, error) {
	key := c.leaderboardKey(period, limit, offset)

	// Try cache
	cached, err := c.rdb.Get(ctx, key).Bytes()
	if err == nil {
		var result domain.LeaderboardResult
		if err := json.Unmarshal(cached, &result); err == nil {
			return &result, nil
		}
		c.logger.Warn("failed to unmarshal cached leaderboard", zap.Error(err))
	}

	// Cache miss — fetch from service
	result, err := c.inner.GetLeaderboard(ctx, period, limit, offset)
	if err != nil {
		return nil, err
	}

	// Store in cache (best-effort)
	if data, err := json.Marshal(result); err == nil {
		if err := c.rdb.Set(ctx, key, data, leaderboardTTL).Err(); err != nil {
			c.logger.Warn("failed to cache leaderboard", zap.Error(err))
		}
	}

	return result, nil
}

// GetUserSummary attempts to serve from cache, falling back to the inner service.
func (c *CachedLeaderboardService) GetUserSummary(ctx context.Context, userID uuid.UUID, period domain.LeaderboardPeriod) (*domain.UserPointsSummary, error) {
	key := c.userSummaryKey(userID, period)

	// Try cache
	cached, err := c.rdb.Get(ctx, key).Bytes()
	if err == nil {
		var summary domain.UserPointsSummary
		if err := json.Unmarshal(cached, &summary); err == nil {
			return &summary, nil
		}
		c.logger.Warn("failed to unmarshal cached user summary", zap.Error(err))
	}

	// Cache miss
	summary, err := c.inner.GetUserSummary(ctx, userID, period)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(summary); err == nil {
		if err := c.rdb.Set(ctx, key, data, userSummaryTTL).Err(); err != nil {
			c.logger.Warn("failed to cache user summary", zap.Error(err))
		}
	}

	return summary, nil
}

// InvalidateLeaderboard clears all cached leaderboard data.
// Called after a refresh cycle completes.
func (c *CachedLeaderboardService) InvalidateLeaderboard(ctx context.Context) {
	pattern := "leaderboard:*"
	iter := c.rdb.Scan(ctx, 0, pattern, 500).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		c.rdb.Del(ctx, keys...)
	}
}

func (c *CachedLeaderboardService) leaderboardKey(period domain.LeaderboardPeriod, limit, offset int) string {
	return fmt.Sprintf("leaderboard:%s:%d:%d", period, limit, offset)
}

func (c *CachedLeaderboardService) userSummaryKey(userID uuid.UUID, period domain.LeaderboardPeriod) string {
	return fmt.Sprintf("leaderboard:user:%s:%s", userID.String(), period)
}
