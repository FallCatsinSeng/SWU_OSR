package repository

import (
	"context"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationRepo implements domain.NotificationRepository using pgxpool.
type NotificationRepo struct {
	pool *pgxpool.Pool
}

// NewNotificationRepo creates a new notification repository.
func NewNotificationRepo(pool *pgxpool.Pool) domain.NotificationRepository {
	return &NotificationRepo{pool: pool}
}

// Create inserts a new notification.
func (r *NotificationRepo) Create(ctx context.Context, notif *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, reference_id, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		notif.ID, notif.UserID, notif.Type, notif.ReferenceID,
		notif.Message, notif.IsRead, notif.CreatedAt,
	)
	return err
}

// ListByUserID retrieves notifications for a user ordered by created_at DESC.
func (r *NotificationRepo) ListByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	query := `
		SELECT id, user_id, type, reference_id, message, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.ReferenceID,
			&n.Message, &n.IsRead, &n.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	if notifications == nil {
		notifications = []domain.Notification{}
	}
	return notifications, rows.Err()
}

// MarkRead marks a notification as read for the given user.
func (r *NotificationRepo) MarkRead(ctx context.Context, userID uuid.UUID, notifID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	_, err := r.pool.Exec(ctx, query, notifID, userID)
	return err
}
