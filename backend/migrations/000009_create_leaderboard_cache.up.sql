-- Materialized view-style cache table for leaderboard points.
-- This table is periodically refreshed by the application to avoid
-- expensive real-time aggregation on every leaderboard request.

CREATE TABLE IF NOT EXISTS leaderboard_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    period VARCHAR(20) NOT NULL,          -- 'weekly', 'monthly', 'all_time', 'semester'
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    push_points INT NOT NULL DEFAULT 0,
    pr_points INT NOT NULL DEFAULT 0,
    forum_points INT NOT NULL DEFAULT 0,
    other_points INT NOT NULL DEFAULT 0,  -- showcase additions, streak bonuses
    total_points INT NOT NULL DEFAULT 0,
    streak_days INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_user_period UNIQUE (user_id, period, period_start)
);

-- Index for efficient leaderboard queries (ranked by total_points desc).
CREATE INDEX idx_leaderboard_period_points 
    ON leaderboard_points (period, period_start, total_points DESC);

-- Index for user-specific lookups.
CREATE INDEX idx_leaderboard_user 
    ON leaderboard_points (user_id, period, period_start);
