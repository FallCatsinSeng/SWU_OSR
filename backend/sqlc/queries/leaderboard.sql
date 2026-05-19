-- name: GetLeaderboard :many
-- Retrieves the leaderboard for a specific period window, ranked by total_points.
SELECT 
    lp.user_id,
    u.alias,
    u.avatar_url,
    lp.push_points,
    lp.pr_points,
    lp.forum_points,
    lp.other_points,
    lp.total_points,
    lp.streak_days
FROM leaderboard_points lp
JOIN users u ON lp.user_id = u.id
WHERE lp.period = $1 
  AND lp.period_start = $2
ORDER BY lp.total_points DESC, u.alias ASC
LIMIT $3 OFFSET $4;

-- name: GetUserPoints :one
-- Retrieves a single user's points for a specific period window.
SELECT 
    lp.user_id,
    lp.push_points,
    lp.pr_points,
    lp.forum_points,
    lp.other_points,
    lp.total_points,
    lp.streak_days
FROM leaderboard_points lp
WHERE lp.user_id = $1 
  AND lp.period = $2 
  AND lp.period_start = $3;

-- name: GetUserRank :one
-- Returns the rank of a user within a given period.
SELECT COUNT(*) + 1 AS rank
FROM leaderboard_points
WHERE period = $1 
  AND period_start = $2 
  AND total_points > (
    SELECT COALESCE(total_points, 0) 
    FROM leaderboard_points 
    WHERE user_id = $3 AND period = $1 AND period_start = $2
  );

-- name: UpsertLeaderboardPoints :exec
-- Inserts or updates leaderboard points for a user in a given period.
INSERT INTO leaderboard_points (
    user_id, period, period_start, period_end,
    push_points, pr_points, forum_points, other_points,
    total_points, streak_days, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
ON CONFLICT (user_id, period, period_start)
DO UPDATE SET
    period_end = EXCLUDED.period_end,
    push_points = EXCLUDED.push_points,
    pr_points = EXCLUDED.pr_points,
    forum_points = EXCLUDED.forum_points,
    other_points = EXCLUDED.other_points,
    total_points = EXCLUDED.total_points,
    streak_days = EXCLUDED.streak_days,
    updated_at = NOW();

-- name: CountUserPushEvents :one
-- Counts push events for a user within a time range.
SELECT COUNT(*) FROM activity_logs
WHERE user_id = $1 AND event_type = 'push' AND created_at >= $2 AND created_at < $3;

-- name: CountUserPREvents :one
-- Counts pull_request events for a user within a time range.
SELECT COUNT(*) FROM activity_logs
WHERE user_id = $1 AND event_type = 'pull_request' AND created_at >= $2 AND created_at < $3;

-- name: CountUserThreads :one
-- Counts threads created by a user within a time range.
SELECT COUNT(*) FROM threads
WHERE author_id = $1 AND created_at >= $2 AND created_at < $3 AND deleted_at IS NULL;

-- name: CountUserComments :one
-- Counts comments posted by a user within a time range.
SELECT COUNT(*) FROM comments
WHERE author_id = $1 AND created_at >= $2 AND created_at < $3 AND deleted_at IS NULL;

-- name: CountUserShowcaseRepos :one
-- Counts showcase repos added by a user within a time range.
SELECT COUNT(*) FROM showcase_repos
WHERE user_id = $1 AND created_at >= $2 AND created_at < $3 AND deleted_at IS NULL;

-- name: GetUserActiveDays :many
-- Returns distinct active days for a user (for streak calculation).
SELECT DISTINCT DATE(created_at) AS active_date
FROM activity_logs
WHERE user_id = $1 AND created_at >= $2
ORDER BY active_date DESC;

-- name: GetAllActiveUserIDs :many
-- Returns all user IDs that have any activity in the given time range.
SELECT DISTINCT user_id FROM activity_logs
WHERE created_at >= $1 AND created_at < $2
UNION
SELECT DISTINCT author_id FROM threads
WHERE created_at >= $1 AND created_at < $2 AND deleted_at IS NULL
UNION
SELECT DISTINCT author_id FROM comments
WHERE created_at >= $1 AND created_at < $2 AND deleted_at IS NULL
UNION
SELECT DISTINCT user_id FROM showcase_repos
WHERE created_at >= $1 AND created_at < $2 AND deleted_at IS NULL;
