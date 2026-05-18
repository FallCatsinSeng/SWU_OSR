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

// ThreadRepo implements domain.ThreadRepository using pgxpool.
type ThreadRepo struct {
	pool *pgxpool.Pool
}

// NewThreadRepo creates a new thread repository.
func NewThreadRepo(pool *pgxpool.Pool) domain.ThreadRepository {
	return &ThreadRepo{pool: pool}
}

// Create inserts a new thread.
func (r *ThreadRepo) Create(ctx context.Context, thread *domain.Thread) error {
	query := `
		INSERT INTO threads (id, showcase_repo_id, author_id, title, body, comment_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.pool.Exec(ctx, query,
		thread.ID, thread.ShowcaseRepoID, thread.AuthorID, thread.Title,
		thread.Body, thread.CommentCount, thread.CreatedAt, thread.UpdatedAt,
	)
	return err
}

// GetByID retrieves a thread by its ID with author info.
func (r *ThreadRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Thread, error) {
	query := `
		SELECT t.id, t.showcase_repo_id, t.author_id, t.title, t.body, t.comment_count,
			t.created_at, t.updated_at
		FROM threads t
		WHERE t.id = $1 AND t.deleted_at IS NULL`

	var thread domain.Thread
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&thread.ID, &thread.ShowcaseRepoID, &thread.AuthorID, &thread.Title,
		&thread.Body, &thread.CommentCount, &thread.CreatedAt, &thread.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &thread, nil
}

// ListByRepoID returns paginated threads for a showcase repo ordered by created_at DESC.
func (r *ThreadRepo) ListByRepoID(ctx context.Context, repoID uuid.UUID, cursor time.Time, limit int) ([]domain.Thread, error) {
	query := `
		SELECT t.id, t.showcase_repo_id, t.author_id, t.title, t.body, t.comment_count,
			t.created_at, t.updated_at
		FROM threads t
		WHERE t.showcase_repo_id = $1 AND t.created_at < $2 AND t.deleted_at IS NULL
		ORDER BY t.created_at DESC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, repoID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []domain.Thread
	for rows.Next() {
		var thread domain.Thread
		if err := rows.Scan(
			&thread.ID, &thread.ShowcaseRepoID, &thread.AuthorID, &thread.Title,
			&thread.Body, &thread.CommentCount, &thread.CreatedAt, &thread.UpdatedAt,
		); err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	if threads == nil {
		threads = []domain.Thread{}
	}
	return threads, rows.Err()
}

// IncrementCommentCount increments the comment count for a thread.
func (r *ThreadRepo) IncrementCommentCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE threads SET comment_count = comment_count + 1 WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
