package repository

import (
	"context"
	"errors"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ShowcaseRepo implements domain.ShowcaseRepository using pgxpool.
type ShowcaseRepo struct {
	pool *pgxpool.Pool
}

// NewShowcaseRepo creates a new showcase repository.
func NewShowcaseRepo(pool *pgxpool.Pool) domain.ShowcaseRepository {
	return &ShowcaseRepo{pool: pool}
}

// Create inserts a new showcase repo.
func (r *ShowcaseRepo) Create(ctx context.Context, repo *domain.ShowcaseRepo) error {
	query := `
		INSERT INTO showcase_repos (id, user_id, github_repo_id, repo_name, repo_full_name,
			description, language, html_url, academic_tag, webhook_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.pool.Exec(ctx, query,
		repo.ID, repo.UserID, repo.GitHubRepoID, repo.RepoName, repo.RepoFullName,
		repo.Description, repo.Language, repo.HTMLURL, string(repo.AcademicTag),
		repo.WebhookID, repo.CreatedAt, repo.UpdatedAt,
	)
	return err
}

// GetByID retrieves a showcase repo by its ID.
func (r *ShowcaseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ShowcaseRepo, error) {
	query := `
		SELECT id, user_id, github_repo_id, repo_name, repo_full_name, description, language,
			html_url, academic_tag, webhook_id, created_at, updated_at
		FROM showcase_repos
		WHERE id = $1 AND deleted_at IS NULL`

	var repo domain.ShowcaseRepo
	var tag string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&repo.ID, &repo.UserID, &repo.GitHubRepoID, &repo.RepoName, &repo.RepoFullName,
		&repo.Description, &repo.Language, &repo.HTMLURL, &tag,
		&repo.WebhookID, &repo.CreatedAt, &repo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	repo.AcademicTag = domain.AcademicTag(tag)
	return &repo, nil
}

// GetByUserID retrieves all active showcase repos for a user.
func (r *ShowcaseRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.ShowcaseRepo, error) {
	query := `
		SELECT id, user_id, github_repo_id, repo_name, repo_full_name, description, language,
			html_url, academic_tag, webhook_id, created_at, updated_at
		FROM showcase_repos
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []domain.ShowcaseRepo
	for rows.Next() {
		var repo domain.ShowcaseRepo
		var tag string
		if err := rows.Scan(
			&repo.ID, &repo.UserID, &repo.GitHubRepoID, &repo.RepoName, &repo.RepoFullName,
			&repo.Description, &repo.Language, &repo.HTMLURL, &tag,
			&repo.WebhookID, &repo.CreatedAt, &repo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		repo.AcademicTag = domain.AcademicTag(tag)
		repos = append(repos, repo)
	}

	if repos == nil {
		repos = []domain.ShowcaseRepo{}
	}
	return repos, rows.Err()
}

// GetByUserAndRepoFullName retrieves a showcase repo by user and full name.
func (r *ShowcaseRepo) GetByUserAndRepoFullName(ctx context.Context, userID uuid.UUID, repoFullName string) (*domain.ShowcaseRepo, error) {
	query := `
		SELECT id, user_id, github_repo_id, repo_name, repo_full_name, description, language,
			html_url, academic_tag, webhook_id, created_at, updated_at
		FROM showcase_repos
		WHERE user_id = $1 AND repo_full_name = $2 AND deleted_at IS NULL`

	var repo domain.ShowcaseRepo
	var tag string
	err := r.pool.QueryRow(ctx, query, userID, repoFullName).Scan(
		&repo.ID, &repo.UserID, &repo.GitHubRepoID, &repo.RepoName, &repo.RepoFullName,
		&repo.Description, &repo.Language, &repo.HTMLURL, &tag,
		&repo.WebhookID, &repo.CreatedAt, &repo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	repo.AcademicTag = domain.AcademicTag(tag)
	return &repo, nil
}

// SoftDelete soft-deletes a showcase repo by ID.
func (r *ShowcaseRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE showcase_repos SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// SoftDeleteByUser soft-deletes a specific showcase repo for a user.
func (r *ShowcaseRepo) SoftDeleteByUser(ctx context.Context, userID uuid.UUID, repoID uuid.UUID) error {
	query := `UPDATE showcase_repos SET deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, repoID, userID)
	return err
}


// GetByUserAndRepoFullNameIncludeDeleted retrieves a soft-deleted showcase repo by user and full name OR github_repo_id.
func (r *ShowcaseRepo) GetByUserAndRepoFullNameIncludeDeleted(ctx context.Context, userID uuid.UUID, repoFullName string) (*domain.ShowcaseRepo, error) {
	query := `
		SELECT id, user_id, github_repo_id, repo_name, repo_full_name, description, language,
			html_url, academic_tag, webhook_id, created_at, updated_at
		FROM showcase_repos
		WHERE user_id = $1 AND repo_full_name = $2 AND deleted_at IS NOT NULL
		ORDER BY deleted_at DESC
		LIMIT 1`

	var repo domain.ShowcaseRepo
	var tag string
	err := r.pool.QueryRow(ctx, query, userID, repoFullName).Scan(
		&repo.ID, &repo.UserID, &repo.GitHubRepoID, &repo.RepoName, &repo.RepoFullName,
		&repo.Description, &repo.Language, &repo.HTMLURL, &tag,
		&repo.WebhookID, &repo.CreatedAt, &repo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	repo.AcademicTag = domain.AcademicTag(tag)
	return &repo, nil
}

// Restore un-deletes a soft-deleted showcase repo and updates its fields.
func (r *ShowcaseRepo) Restore(ctx context.Context, repo *domain.ShowcaseRepo) error {
	query := `
		UPDATE showcase_repos
		SET deleted_at = NULL, academic_tag = $2, description = $3, language = $4,
			html_url = $5, webhook_id = $6, updated_at = $7, repo_name = $8, repo_full_name = $9
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		repo.ID, string(repo.AcademicTag), repo.Description, repo.Language,
		repo.HTMLURL, repo.WebhookID, repo.UpdatedAt, repo.RepoName, repo.RepoFullName,
	)
	return err
}

// GetByUserAndGitHubRepoIDIncludeDeleted retrieves a soft-deleted showcase repo by github_repo_id.
func (r *ShowcaseRepo) GetByUserAndGitHubRepoIDIncludeDeleted(ctx context.Context, userID uuid.UUID, githubRepoID int64) (*domain.ShowcaseRepo, error) {
	query := `
		SELECT id, user_id, github_repo_id, repo_name, repo_full_name, description, language,
			html_url, academic_tag, webhook_id, created_at, updated_at
		FROM showcase_repos
		WHERE user_id = $1 AND github_repo_id = $2 AND deleted_at IS NOT NULL
		ORDER BY deleted_at DESC
		LIMIT 1`

	var repo domain.ShowcaseRepo
	var tag string
	err := r.pool.QueryRow(ctx, query, userID, githubRepoID).Scan(
		&repo.ID, &repo.UserID, &repo.GitHubRepoID, &repo.RepoName, &repo.RepoFullName,
		&repo.Description, &repo.Language, &repo.HTMLURL, &tag,
		&repo.WebhookID, &repo.CreatedAt, &repo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	repo.AcademicTag = domain.AcademicTag(tag)
	return &repo, nil
}
