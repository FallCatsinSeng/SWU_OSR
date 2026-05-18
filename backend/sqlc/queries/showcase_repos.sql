-- name: CreateShowcaseRepo :one
INSERT INTO showcase_repos (user_id, github_repo_id, repo_name, repo_full_name, description, language, html_url, academic_tag, webhook_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetShowcaseReposByUserID :many
SELECT * FROM showcase_repos WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: GetShowcaseRepoByUserAndFullName :one
SELECT * FROM showcase_repos WHERE user_id = $1 AND repo_full_name = $2 AND deleted_at IS NULL;

-- name: SoftDeleteShowcaseRepo :exec
UPDATE showcase_repos SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: SoftDeleteShowcaseRepoByUser :exec
UPDATE showcase_repos SET deleted_at = NOW(), updated_at = NOW() WHERE user_id = $1 AND id = $2;
