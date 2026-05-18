-- name: CreateNotification :one
INSERT INTO notifications (user_id, type, reference_id, message)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListNotificationsByUserID :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: MarkNotificationRead :exec
UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2;
