-- name: CreateThread :one
INSERT INTO threads (showcase_repo_id, author_id, title, body)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetThreadByID :one
SELECT * FROM threads WHERE id = $1 AND deleted_at IS NULL;

-- name: ListThreadsByRepoID :many
SELECT * FROM threads
WHERE showcase_repo_id = $1 AND deleted_at IS NULL AND created_at < $2
ORDER BY created_at DESC
LIMIT $3;

-- name: IncrementThreadCommentCount :exec
UPDATE threads SET comment_count = comment_count + 1, updated_at = NOW() WHERE id = $1;
