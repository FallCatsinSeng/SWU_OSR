-- name: InsertActivityLog :one
INSERT INTO activity_logs (user_id, showcase_repo_id, event_type, summary, metadata, github_event_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetActivityFeed :many
SELECT al.*, u.alias AS user_alias, u.avatar_url
FROM activity_logs al
JOIN users u ON u.id = al.user_id
WHERE al.created_at < $1
ORDER BY al.created_at DESC
LIMIT $2;

-- name: GetUserActivityFeed :many
SELECT al.*, u.alias AS user_alias, u.avatar_url
FROM activity_logs al
JOIN users u ON u.id = al.user_id
WHERE al.user_id = $1 AND al.created_at < $2
ORDER BY al.created_at DESC
LIMIT $3;

-- name: GetActivityLogByGitHubEventID :one
SELECT * FROM activity_logs WHERE github_event_id = $1;
