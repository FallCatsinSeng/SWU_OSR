package service

import (
	"context"
	"testing"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// mockLeaderboardRepo implements domain.LeaderboardRepository for testing.
type mockLeaderboardRepo struct {
	entries         []domain.LeaderboardEntry
	userPoints      map[uuid.UUID]*domain.UserPointsSummary
	userStreaks      map[uuid.UUID]int
	upsertedPoints  map[uuid.UUID]int
	activeUserIDs   []uuid.UUID
	pushPerRepo     map[uuid.UUID][]domain.RepoEventCount
	prPerRepo       map[uuid.UUID][]domain.RepoEventCount
	threadCounts    map[uuid.UUID]int
	commentCounts   map[uuid.UUID]int
	showcaseCounts  map[uuid.UUID]int
}

func newMockLeaderboardRepo() *mockLeaderboardRepo {
	return &mockLeaderboardRepo{
		entries:        []domain.LeaderboardEntry{},
		userPoints:     make(map[uuid.UUID]*domain.UserPointsSummary),
		userStreaks:    make(map[uuid.UUID]int),
		upsertedPoints: make(map[uuid.UUID]int),
		pushPerRepo:    make(map[uuid.UUID][]domain.RepoEventCount),
		prPerRepo:      make(map[uuid.UUID][]domain.RepoEventCount),
		threadCounts:   make(map[uuid.UUID]int),
		commentCounts:  make(map[uuid.UUID]int),
		showcaseCounts: make(map[uuid.UUID]int),
	}
}

func (m *mockLeaderboardRepo) GetLeaderboard(_ context.Context, _, _ time.Time, limit, offset int) ([]domain.LeaderboardEntry, error) {
	end := offset + limit
	if end > len(m.entries) {
		end = len(m.entries)
	}
	if offset > len(m.entries) {
		return []domain.LeaderboardEntry{}, nil
	}
	return m.entries[offset:end], nil
}

func (m *mockLeaderboardRepo) GetUserPoints(_ context.Context, userID uuid.UUID, _, _ time.Time) (*domain.UserPointsSummary, error) {
	if summary, ok := m.userPoints[userID]; ok {
		return summary, nil
	}
	return &domain.UserPointsSummary{UserID: userID}, nil
}

func (m *mockLeaderboardRepo) GetUserStreak(_ context.Context, userID uuid.UUID) (int, error) {
	return m.userStreaks[userID], nil
}

func (m *mockLeaderboardRepo) CountUserPushEvents(_ context.Context, userID uuid.UUID, _, _ time.Time) (int, error) {
	return 0, nil
}

func (m *mockLeaderboardRepo) CountUserPREvents(_ context.Context, userID uuid.UUID, _, _ time.Time) (int, error) {
	return 0, nil
}

func (m *mockLeaderboardRepo) CountUserThreads(_ context.Context, userID uuid.UUID, _, _ time.Time) (int, error) {
	return m.threadCounts[userID], nil
}

func (m *mockLeaderboardRepo) CountUserComments(_ context.Context, userID uuid.UUID, _, _ time.Time) (int, error) {
	return m.commentCounts[userID], nil
}

func (m *mockLeaderboardRepo) CountUserShowcaseRepos(_ context.Context, userID uuid.UUID, _, _ time.Time) (int, error) {
	return m.showcaseCounts[userID], nil
}

func (m *mockLeaderboardRepo) CountUserPushEventsPerRepo(_ context.Context, userID uuid.UUID, _, _ time.Time) ([]domain.RepoEventCount, error) {
	return m.pushPerRepo[userID], nil
}

func (m *mockLeaderboardRepo) CountUserPREventsPerRepo(_ context.Context, userID uuid.UUID, _, _ time.Time) ([]domain.RepoEventCount, error) {
	return m.prPerRepo[userID], nil
}

func (m *mockLeaderboardRepo) GetAllActiveUserIDs(_ context.Context, _, _ time.Time) ([]uuid.UUID, error) {
	return m.activeUserIDs, nil
}

func (m *mockLeaderboardRepo) UpsertPoints(_ context.Context, userID uuid.UUID, _ domain.LeaderboardPeriod, _, _ time.Time, _, _, _, _, totalPts, _ int) error {
	m.upsertedPoints[userID] = totalPts
	return nil
}

func setupLeaderboardService() (*leaderboardService, *mockLeaderboardRepo) {
	repo := newMockLeaderboardRepo()
	svc := NewLeaderboardService(repo, zap.NewNop())
	return svc, repo
}

// --- GetLeaderboard ---

func TestGetLeaderboard_ReturnsPaginatedEntries(t *testing.T) {
	svc, repo := setupLeaderboardService()

	userID1 := uuid.New()
	userID2 := uuid.New()
	repo.entries = []domain.LeaderboardEntry{
		{Rank: 1, UserID: userID1, Alias: "alice", TotalPoints: 100},
		{Rank: 2, UserID: userID2, Alias: "bob", TotalPoints: 80},
	}

	result, err := svc.GetLeaderboard(context.Background(), domain.PeriodWeekly, 10, 0)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, domain.PeriodWeekly, result.Period)
	assert.Len(t, result.Entries, 2)
	assert.Equal(t, "alice", result.Entries[0].Alias)
}

func TestGetLeaderboard_RespectsLimitAndOffset(t *testing.T) {
	svc, repo := setupLeaderboardService()

	for i := 0; i < 5; i++ {
		repo.entries = append(repo.entries, domain.LeaderboardEntry{
			Rank:        i + 1,
			UserID:      uuid.New(),
			TotalPoints: (5 - i) * 10,
		})
	}

	result, err := svc.GetLeaderboard(context.Background(), domain.PeriodMonthly, 2, 1)

	require.NoError(t, err)
	assert.Len(t, result.Entries, 2)
}

func TestGetLeaderboard_EmptyWhenNoEntries(t *testing.T) {
	svc, _ := setupLeaderboardService()

	result, err := svc.GetLeaderboard(context.Background(), domain.PeriodAllTime, 20, 0)

	require.NoError(t, err)
	assert.Empty(t, result.Entries)
}

// --- GetUserSummary ---

func TestGetUserSummary_ReturnsPointsWithStreak(t *testing.T) {
	svc, repo := setupLeaderboardService()

	userID := uuid.New()
	repo.userPoints[userID] = &domain.UserPointsSummary{
		UserID:      userID,
		TotalPoints: 45,
	}
	repo.userStreaks[userID] = 10

	summary, err := svc.GetUserSummary(context.Background(), userID, domain.PeriodWeekly)

	require.NoError(t, err)
	require.NotNil(t, summary)
	assert.Equal(t, 45, summary.TotalPoints)
	assert.Equal(t, 10, summary.StreakDays)
}

func TestGetUserSummary_AllPeriods(t *testing.T) {
	svc, _ := setupLeaderboardService()
	userID := uuid.New()

	periods := []domain.LeaderboardPeriod{
		domain.PeriodWeekly,
		domain.PeriodMonthly,
		domain.PeriodSemester,
		domain.PeriodAllTime,
	}

	for _, period := range periods {
		result, err := svc.GetUserSummary(context.Background(), userID, period)
		require.NoError(t, err, "period: %s", period)
		require.NotNil(t, result, "period: %s", period)
	}
}

// --- periodWindow ---

func TestPeriodWindow_Weekly_StartsOnMonday(t *testing.T) {
	svc, _ := setupLeaderboardService()

	from, to := svc.periodWindow(domain.PeriodWeekly)

	// The "from" date should always be a Monday
	assert.Equal(t, time.Monday, from.Weekday())
	// The window should be exactly 7 days
	assert.Equal(t, 7*24*time.Hour, to.Sub(from))
}

func TestPeriodWindow_Monthly_StartsOnFirstDay(t *testing.T) {
	svc, _ := setupLeaderboardService()

	from, _ := svc.periodWindow(domain.PeriodMonthly)

	assert.Equal(t, 1, from.Day())
}

func TestPeriodWindow_AllTime_UsesFixedSentinelDates(t *testing.T) {
	svc, _ := setupLeaderboardService()

	from, to := svc.periodWindow(domain.PeriodAllTime)

	// Should use fixed sentinel dates so reads/writes always match
	assert.Equal(t, 2020, from.Year())
	assert.Equal(t, 2099, to.Year())
}

func TestPeriodWindow_DefaultsToWeekly(t *testing.T) {
	svc, _ := setupLeaderboardService()

	fromDefault, toDefault := svc.periodWindow("")
	fromWeekly, toWeekly := svc.periodWindow(domain.PeriodWeekly)

	assert.Equal(t, fromWeekly, fromDefault)
	assert.Equal(t, toWeekly, toDefault)
}

// --- applyCappedCount (anti-gaming logic) ---

func TestApplyCappedCount_CapsPerRepoPushCount(t *testing.T) {
	svc, _ := setupLeaderboardService()

	now := time.Now()
	from := now.AddDate(0, 0, -6) // weekly window
	to := now

	// User pushed 50 times to one repo (cap is 15/week)
	perRepo := []domain.RepoEventCount{
		{RepoID: uuid.New(), Count: 50},
	}

	total := svc.applyCappedCount(perRepo, domain.MaxPushPerRepoPerWeek, from, to)

	// Should be capped at MaxPushPerRepoPerWeek (15)
	assert.Equal(t, domain.MaxPushPerRepoPerWeek, total)
}

func TestApplyCappedCount_SumsMultipleRepos(t *testing.T) {
	svc, _ := setupLeaderboardService()

	now := time.Now()
	from := now.AddDate(0, 0, -6)
	to := now

	perRepo := []domain.RepoEventCount{
		{RepoID: uuid.New(), Count: 5},  // under cap
		{RepoID: uuid.New(), Count: 3},  // under cap
		{RepoID: uuid.New(), Count: 20}, // over cap (15)
	}

	total := svc.applyCappedCount(perRepo, domain.MaxPushPerRepoPerWeek, from, to)

	// 5 + 3 + 15 (capped) = 23
	assert.Equal(t, 23, total)
}

// --- capDailyPoints ---

func TestCapDailyPoints_CapsAtMaxDailyLimit(t *testing.T) {
	svc, _ := setupLeaderboardService()

	now := time.Now()
	from := now.AddDate(0, 0, -6) // 7-day window
	to := now

	// 1000 push events, each worth 3 points, but daily cap is 30 pts
	// Max events/day = 30/3 = 10 events/day
	// Total max = 10 * 7 days = 70 events
	// Points = 70 * 3 = 210
	points := svc.capDailyPoints(1000, domain.PointsPush, from, to)

	assert.Equal(t, 210, points)
}

func TestCapDailyPoints_AllowsUnderLimit(t *testing.T) {
	svc, _ := setupLeaderboardService()

	now := time.Now()
	from := now.AddDate(0, 0, -6)
	to := now

	// Only 3 pushes at 3 pts each = 9 pts total, well under daily cap
	points := svc.capDailyPoints(3, domain.PointsPush, from, to)

	assert.Equal(t, 9, points)
}
