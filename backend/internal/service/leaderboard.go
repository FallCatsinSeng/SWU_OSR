package service

import (
	"context"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// leaderboardService implements domain.LeaderboardService.
type leaderboardService struct {
	repo   domain.LeaderboardRepository
	logger *zap.Logger
}

// NewLeaderboardService creates a new leaderboard service.
func NewLeaderboardService(repo domain.LeaderboardRepository, logger *zap.Logger) domain.LeaderboardService {
	return &leaderboardService{
		repo:   repo,
		logger: logger,
	}
}

// GetLeaderboard returns the leaderboard for a given period.
func (s *leaderboardService) GetLeaderboard(ctx context.Context, period domain.LeaderboardPeriod, limit, offset int) (*domain.LeaderboardResult, error) {
	from, to := s.periodWindow(period)

	entries, err := s.repo.GetLeaderboard(ctx, from, to, limit, offset)
	if err != nil {
		s.logger.Error("failed to get leaderboard", zap.Error(err))
		return nil, err
	}

	return &domain.LeaderboardResult{
		Period:  period,
		From:    from,
		To:      to,
		Entries: entries,
	}, nil
}

// GetUserSummary returns a user's personal points summary for a given period.
func (s *leaderboardService) GetUserSummary(ctx context.Context, userID uuid.UUID, period domain.LeaderboardPeriod) (*domain.UserPointsSummary, error) {
	from, to := s.periodWindow(period)

	summary, err := s.repo.GetUserPoints(ctx, userID, from, to)
	if err != nil {
		s.logger.Error("failed to get user points", zap.Error(err), zap.String("user_id", userID.String()))
		return nil, err
	}

	// Get live streak
	streak, err := s.repo.GetUserStreak(ctx, userID)
	if err != nil {
		s.logger.Warn("failed to get user streak", zap.Error(err))
		// Non-fatal: continue with cached streak
	} else {
		summary.StreakDays = streak
	}

	return summary, nil
}

// RefreshLeaderboard recomputes points for all active users in the given period.
// This should be called periodically (e.g., every 15 minutes via cron/scheduler).
func (s *leaderboardService) RefreshLeaderboard(ctx context.Context, period domain.LeaderboardPeriod) error {
	from, to := s.periodWindow(period)

	// Get all users with activity in this period
	userIDs, err := s.repo.GetAllActiveUserIDs(ctx, from, to)
	if err != nil {
		s.logger.Error("failed to get active user IDs", zap.Error(err))
		return err
	}

	s.logger.Info("refreshing leaderboard",
		zap.String("period", string(period)),
		zap.Int("active_users", len(userIDs)),
	)

	for _, userID := range userIDs {
		if err := s.computeAndStoreUserPoints(ctx, userID, period, from, to); err != nil {
			s.logger.Error("failed to compute points for user",
				zap.Error(err),
				zap.String("user_id", userID.String()),
			)
			// Continue processing other users
			continue
		}
	}

	return nil
}

// computeAndStoreUserPoints calculates and persists a user's points.
func (s *leaderboardService) computeAndStoreUserPoints(ctx context.Context, userID uuid.UUID, period domain.LeaderboardPeriod, from, to time.Time) error {
	// Count activities
	pushCount, err := s.repo.CountUserPushEvents(ctx, userID, from, to)
	if err != nil {
		return err
	}

	prCount, err := s.repo.CountUserPREvents(ctx, userID, from, to)
	if err != nil {
		return err
	}

	threadCount, err := s.repo.CountUserThreads(ctx, userID, from, to)
	if err != nil {
		return err
	}

	commentCount, err := s.repo.CountUserComments(ctx, userID, from, to)
	if err != nil {
		return err
	}

	showcaseCount, err := s.repo.CountUserShowcaseRepos(ctx, userID, from, to)
	if err != nil {
		return err
	}

	streak, err := s.repo.GetUserStreak(ctx, userID)
	if err != nil {
		streak = 0
	}

	// Calculate points with daily cap applied
	pushPts := s.capDailyPoints(pushCount, domain.PointsPush, from, to)
	prPts := prCount * domain.PointsPROpened
	forumPts := (threadCount * domain.PointsThreadCreated) + (commentCount * domain.PointsCommentPosted)
	otherPts := showcaseCount * domain.PointsShowcaseAdded

	// Add streak bonus if threshold met
	if streak >= domain.StreakThresholdDays {
		streakBonuses := streak / domain.StreakThresholdDays
		otherPts += streakBonuses * domain.PointsStreakBonus
	}

	totalPts := pushPts + prPts + forumPts + otherPts

	// Persist
	return s.repo.UpsertPoints(ctx, userID, period, from, to, pushPts, prPts, forumPts, otherPts, totalPts, streak)
}

// capDailyPoints applies the daily point cap for push events.
// It limits the raw count so that points don't exceed MaxPointsPerDay worth of push points per day.
func (s *leaderboardService) capDailyPoints(count, pointsPerEvent int, from, to time.Time) int {
	days := int(to.Sub(from).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	// Max push events that count per day
	maxEventsPerDay := domain.MaxPointsPerDay / pointsPerEvent
	maxTotalEvents := maxEventsPerDay * days

	if count > maxTotalEvents {
		count = maxTotalEvents
	}

	return count * pointsPerEvent
}

// periodWindow returns the start and end time for a given leaderboard period.
func (s *leaderboardService) periodWindow(period domain.LeaderboardPeriod) (time.Time, time.Time) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	switch period {
	case domain.PeriodWeekly:
		// Start from Monday of the current week
		weekday := int(today.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		from := today.AddDate(0, 0, -(weekday - 1))
		to := from.AddDate(0, 0, 7)
		return from, to

	case domain.PeriodMonthly:
		from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		to := from.AddDate(0, 1, 0)
		return from, to

	case domain.PeriodSemester:
		// Academic semesters: Jan-Jun (even), Jul-Dec (odd)
		if now.Month() <= 6 {
			from := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			to := time.Date(now.Year(), 7, 1, 0, 0, 0, 0, now.Location())
			return from, to
		}
		from := time.Date(now.Year(), 7, 1, 0, 0, 0, 0, now.Location())
		to := time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
		return from, to

	case domain.PeriodAllTime:
		// Use a very early date as start
		from := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		to := now.AddDate(0, 0, 1) // tomorrow
		return from, to

	default:
		// Default to weekly
		weekday := int(today.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		from := today.AddDate(0, 0, -(weekday - 1))
		to := from.AddDate(0, 0, 7)
		return from, to
	}
}
