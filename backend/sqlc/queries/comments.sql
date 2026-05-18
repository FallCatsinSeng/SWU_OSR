-- name: CreateComment :one
INSERT INTO comments (thread_id, author_id, parent_id, body)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCommentsByThreadID :many
SELECT * FROM comments
WHERE thread_id = $1 AND deleted_at IS NULL AND created_at > $2
ORDER BY created_at ASC
LIMIT $3;
