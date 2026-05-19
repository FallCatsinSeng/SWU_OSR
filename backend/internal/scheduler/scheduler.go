package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"go.uber.org/zap"
)

// LeaderboardRefresher defines the method needed by the scheduler.
type LeaderboardRefresher interface {
	RefreshLeaderboard(ctx context.Context, period domain.LeaderboardPeriod) error
}

// Scheduler manages periodic background tasks.
type Scheduler struct {
	refresher LeaderboardRefresher
	logger    *zap.Logger
	interval  time.Duration
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

// New creates a new scheduler.
func New(refresher LeaderboardRefresher, logger *zap.Logger, interval time.Duration) *Scheduler {
	return &Scheduler{
		refresher: refresher,
		logger:    logger,
		interval:  interval,
		stopCh:    make(chan struct{}),
	}
}

// Start begins the periodic leaderboard refresh in a background goroutine.
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go s.run()
	s.logger.Info("leaderboard scheduler started", zap.Duration("interval", s.interval))
}

// Stop gracefully stops the scheduler and waits for the goroutine to finish.
func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.logger.Info("leaderboard scheduler stopped")
}

func (s *Scheduler) run() {
	defer s.wg.Done()

	// Run immediately on startup
	s.refreshAll()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.refreshAll()
		case <-s.stopCh:
			return
		}
	}
}

// refreshAll refreshes all leaderboard periods.
func (s *Scheduler) refreshAll() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	periods := []domain.LeaderboardPeriod{
		domain.PeriodWeekly,
		domain.PeriodMonthly,
		domain.PeriodSemester,
		domain.PeriodAllTime,
	}

	for _, period := range periods {
		start := time.Now()
		if err := s.refresher.RefreshLeaderboard(ctx, period); err != nil {
			s.logger.Error("leaderboard refresh failed",
				zap.String("period", string(period)),
				zap.Error(err),
			)
		} else {
			s.logger.Info("leaderboard refresh completed",
				zap.String("period", string(period)),
				zap.Duration("duration", time.Since(start)),
			)
		}
	}
}
