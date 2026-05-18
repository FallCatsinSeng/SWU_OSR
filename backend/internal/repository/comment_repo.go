package repository

import (
	"context"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CommentRepo implements domain.CommentRepository using pgxpool.
type CommentRepo struct {
	pool *pgxpool.Pool
}

// NewCommentRepo creates a new comment repository.
func NewCommentRepo(pool *pgxpool.Pool) domain.CommentRepository {
	return &CommentRepo{pool: pool}
}

// Create inserts a new comment.
func (r *CommentRepo) Create(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (id, thread_id, author_id, parent_id, body, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		comment.ID, comment.ThreadID, comment.AuthorID, comment.ParentID,
		comment.Body, comment.CreatedAt, comment.UpdatedAt,
	)
	return err
}

// GetByThreadID retrieves comments for a thread ordered by created_at ASC.
func (r *CommentRepo) GetByThreadID(ctx context.Context, threadID uuid.UUID, cursor time.Time, limit int) ([]domain.Comment, error) {
	query := `
		SELECT c.id, c.thread_id, c.author_id, c.parent_id, c.body, c.created_at, c.updated_at
		FROM comments c
		WHERE c.thread_id = $1 AND c.created_at > $2 AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, threadID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID, &comment.ThreadID, &comment.AuthorID, &comment.ParentID,
			&comment.Body, &comment.CreatedAt, &comment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if comments == nil {
		comments = []domain.Comment{}
	}
	return comments, rows.Err()
}
