package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Point values for different activities.
const (
	PointsPush          = 3
	PointsPROpened      = 5
	PointsPRMerged      = 8
	PointsShowcaseAdded = 10
	PointsThreadCreated = 2
	PointsCommentPosted = 1
	PointsStreakBonus   = 15

	// Anti-gaming: maximum points a user can earn per day.
	MaxPointsPerDay = 30

	// Anti-gaming: maximum scoreable push events per repo per week.
	// If someone makes more than this many pushes to a single repo in a week,
	// the excess won't count toward points.
	MaxPushPerRepoPerWeek = 15

	// Anti-gaming: maximum scoreable PR events per repo per week.
	MaxPRPerRepoPerWeek = 5

	// Streak threshold: consecutive days needed for bonus.
	StreakThresholdDays = 7
)

// LeaderboardPeriod represents the time window for a leaderboard.
type LeaderboardPeriod string

const (
	PeriodWeekly   LeaderboardPeriod = "weekly"
	PeriodMonthly  LeaderboardPeriod = "monthly"
	PeriodAllTime  LeaderboardPeriod = "all_time"
	PeriodSemester LeaderboardPeriod = "semester"
)

// LeaderboardEntry represents a single user's position on the leaderboard.
type LeaderboardEntry struct {
	Rank        int       `json:"rank"`
	UserID      uuid.UUID `json:"user_id"`
	Alias       string    `json:"alias"`
	AvatarURL   string    `json:"avatar_url"`
	TotalPoints int       `json:"total_points"`
	PushPoints  int       `json:"push_points"`
	PRPoints    int       `json:"pr_points"`
	ForumPoints int       `json:"forum_points"`
	OtherPoints int       `json:"other_points"`
	StreakDays  int       `json:"streak_days"`
}

// LeaderboardResult is the response for a leaderboard query.
type LeaderboardResult struct {
	Period  LeaderboardPeriod  `json:"period"`
	From    time.Time          `json:"from"`
	To      time.Time          `json:"to"`
	Entries []LeaderboardEntry `json:"entries"`
}

// UserPointsSummary is a user's personal points breakdown.
type UserPointsSummary struct {
	UserID      uuid.UUID `json:"user_id"`
	TotalPoints int       `json:"total_points"`
	PushCount   int       `json:"push_count"`
	PRCount     int       `json:"pr_count"`
	ThreadCount int       `json:"thread_count"`
	CommentCnt  int       `json:"comment_count"`
	ShowcaseCnt int       `json:"showcase_count"`
	StreakDays  int       `json:"streak_days"`
	Rank        int       `json:"rank"`
}

// RepoEventCount holds event counts for a single repository.
type RepoEventCount struct {
	RepoID uuid.UUID
	Count  int
}

// LeaderboardRepository defines data access methods for leaderboard.
type LeaderboardRepository interface {
	// GetLeaderboard retrieves ranked users for a given time window.
	GetLeaderboard(ctx context.Context, from, to time.Time, limit, offset int) ([]LeaderboardEntry, error)

	// GetUserPoints retrieves a single user's points summary for a given time window.
	GetUserPoints(ctx context.Context, userID uuid.UUID, from, to time.Time) (*UserPointsSummary, error)

	// GetUserStreak returns the current consecutive active days for a user.
	GetUserStreak(ctx context.Context, userID uuid.UUID) (int, error)

	// Count methods for computing points from source tables.
	CountUserPushEvents(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error)
	CountUserPREvents(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error)
	CountUserThreads(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error)
	CountUserComments(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error)
	CountUserShowcaseRepos(ctx context.Context, userID uuid.UUID, from, to time.Time) (int, error)

	// Per-repo count methods for anti-gaming caps.
	CountUserPushEventsPerRepo(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]RepoEventCount, error)
	CountUserPREventsPerRepo(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]RepoEventCount, error)

	// GetAllActiveUserIDs returns user IDs with activity in the given time range.
	GetAllActiveUserIDs(ctx context.Context, from, to time.Time) ([]uuid.UUID, error)

	// UpsertPoints inserts or updates computed leaderboard points.
	UpsertPoints(ctx context.Context, userID uuid.UUID, period LeaderboardPeriod, from, to time.Time, pushPts, prPts, forumPts, otherPts, totalPts, streakDays int) error
}

// LeaderboardService defines the business logic for the leaderboard feature.
type LeaderboardService interface {
	// GetLeaderboard returns the leaderboard for a given period.
	GetLeaderboard(ctx context.Context, period LeaderboardPeriod, limit, offset int) (*LeaderboardResult, error)

	// GetUserSummary returns a user's personal points summary for a given period.
	GetUserSummary(ctx context.Context, userID uuid.UUID, period LeaderboardPeriod) (*UserPointsSummary, error)
}
