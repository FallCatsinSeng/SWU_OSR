package repository

import (
	"context"
	"errors"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ActivityRepo implements domain.ActivityRepository using pgxpool.
type ActivityRepo struct {
	pool *pgxpool.Pool
}

// NewActivityRepo creates a new activity repository.
func NewActivityRepo(pool *pgxpool.Pool) domain.ActivityRepository {
	return &ActivityRepo{pool: pool}
}

// Insert inserts a new activity log with deduplication by github_event_id.
func (r *ActivityRepo) Insert(ctx context.Context, log *domain.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (id, user_id, showcase_repo_id, repo_full_name, event_type, summary, metadata, github_event_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (github_event_id) WHERE github_event_id IS NOT NULL DO NOTHING`

	_, err := r.pool.Exec(ctx, query,
		log.ID, log.UserID, log.ShowcaseRepoID, log.RepoFullName, string(log.EventType),
		log.Summary, log.Metadata, log.GitHubEventID, log.CreatedAt,
	)
	return err
}

// GetFeed retrieves a cursor-based paginated activity feed.
func (r *ActivityRepo) GetFeed(ctx context.Context, cursor time.Time, limit int) ([]domain.ActivityItem, error) {
	query := `
		SELECT a.id, a.user_id, u.alias, u.avatar_url, a.event_type,
			a.showcase_repo_id, COALESCE(s.repo_full_name, a.repo_full_name),
			a.repo_full_name, a.summary, a.metadata, a.created_at
		FROM activity_logs a
		JOIN users u ON a.user_id = u.id
		LEFT JOIN showcase_repos s ON a.showcase_repo_id = s.id
		WHERE a.created_at < $1
		ORDER BY a.created_at DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ActivityItem
	for rows.Next() {
		var item domain.ActivityItem
		var eventType string
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserAlias, &item.AvatarURL,
			&eventType, &item.RepoID, &item.RepoName, &item.RepoFullName,
			&item.Summary, &item.Metadata, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.EventType = domain.EventType(eventType)
		items = append(items, item)
	}

	if items == nil {
		items = []domain.ActivityItem{}
	}
	return items, rows.Err()
}

// GetUserFeed retrieves a cursor-based paginated activity feed for a specific user.
func (r *ActivityRepo) GetUserFeed(ctx context.Context, userID uuid.UUID, cursor time.Time, limit int) ([]domain.ActivityItem, error) {
	query := `
		SELECT a.id, a.user_id, u.alias, u.avatar_url, a.event_type,
			a.showcase_repo_id, COALESCE(s.repo_full_name, a.repo_full_name),
			a.repo_full_name, a.summary, a.metadata, a.created_at
		FROM activity_logs a
		JOIN users u ON a.user_id = u.id
		LEFT JOIN showcase_repos s ON a.showcase_repo_id = s.id
		WHERE a.user_id = $1 AND a.created_at < $2
		ORDER BY a.created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, userID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ActivityItem
	for rows.Next() {
		var item domain.ActivityItem
		var eventType string
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserAlias, &item.AvatarURL,
			&eventType, &item.RepoID, &item.RepoName, &item.RepoFullName,
			&item.Summary, &item.Metadata, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.EventType = domain.EventType(eventType)
		items = append(items, item)
	}

	if items == nil {
		items = []domain.ActivityItem{}
	}
	return items, rows.Err()
}

// GetRepoFeed retrieves a cursor-based paginated activity feed for a specific showcase repo.
func (r *ActivityRepo) GetRepoFeed(ctx context.Context, showcaseRepoID uuid.UUID, cursor time.Time, limit int) ([]domain.ActivityItem, error) {
	query := `
		SELECT a.id, a.user_id, u.alias, u.avatar_url, a.event_type,
			a.showcase_repo_id, COALESCE(s.repo_full_name, a.repo_full_name),
			a.repo_full_name, a.summary, a.metadata, a.created_at
		FROM activity_logs a
		JOIN users u ON a.user_id = u.id
		LEFT JOIN showcase_repos s ON a.showcase_repo_id = s.id
		WHERE a.showcase_repo_id = $1 AND a.created_at < $2
		ORDER BY a.created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, showcaseRepoID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ActivityItem
	for rows.Next() {
		var item domain.ActivityItem
		var eventType string
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserAlias, &item.AvatarURL,
			&eventType, &item.RepoID, &item.RepoName, &item.RepoFullName,
			&item.Summary, &item.Metadata, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.EventType = domain.EventType(eventType)
		items = append(items, item)
	}

	if items == nil {
		items = []domain.ActivityItem{}
	}
	return items, rows.Err()
}

// GetByGitHubEventID retrieves an activity log by its GitHub event ID.
func (r *ActivityRepo) GetByGitHubEventID(ctx context.Context, eventID string) (*domain.ActivityLog, error) {
	query := `
		SELECT id, user_id, showcase_repo_id, repo_full_name, event_type, summary, metadata, github_event_id, created_at
		FROM activity_logs
		WHERE github_event_id = $1`

	var log domain.ActivityLog
	var eventType string
	err := r.pool.QueryRow(ctx, query, eventID).Scan(
		&log.ID, &log.UserID, &log.ShowcaseRepoID, &log.RepoFullName, &eventType,
		&log.Summary, &log.Metadata, &log.GitHubEventID, &log.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	log.EventType = domain.EventType(eventType)
	return &log, nil
}
