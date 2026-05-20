-- Performance: Add missing indexes for frequently queried columns.
-- These eliminate sequential scans during leaderboard refresh, activity feeds,
-- community stats, and user profile lookups.

-- Index for GetRepoFeed: activity_logs filtered by showcase_repo_id, ordered by created_at DESC
CREATE INDEX IF NOT EXISTS idx_activity_logs_repo_created
    ON activity_logs (showcase_repo_id, created_at DESC);

-- Index for CountUserPushEvents / CountUserPREvents / GetUserFeed:
-- activity_logs filtered by user_id + event_type, ordered by created_at
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_event_created
    ON activity_logs (user_id, event_type, created_at DESC);

-- Index for CountUserThreads: threads filtered by author_id
CREATE INDEX IF NOT EXISTS idx_threads_author_created
    ON threads (author_id, created_at DESC);

-- Index for CountUserComments: comments filtered by author_id
CREATE INDEX IF NOT EXISTS idx_comments_author_created
    ON comments (author_id, created_at DESC);

-- Index for CountUserShowcaseRepos / GetByUserID: showcase_repos filtered by user_id
CREATE INDEX IF NOT EXISTS idx_showcase_repos_user_active
    ON showcase_repos (user_id, created_at DESC) WHERE deleted_at IS NULL;

-- Index for community stats: active users count
CREATE INDEX IF NOT EXISTS idx_users_active
    ON users (is_active) WHERE deleted_at IS NULL AND is_active = true;

-- Index for PopularRepos correlated subquery: activity count per showcase repo
CREATE INDEX IF NOT EXISTS idx_activity_logs_showcase_repo_id
    ON activity_logs (showcase_repo_id);
