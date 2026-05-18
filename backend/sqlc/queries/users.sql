-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByNIM :one
SELECT * FROM users WHERE nim = $1 AND deleted_at IS NULL;

-- name: GetUserByAlias :one
SELECT * FROM users WHERE alias = $1 AND deleted_at IS NULL;

-- name: GetUserByGitHubUsername :one
SELECT * FROM users WHERE github_username = $1 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (nim, full_name, major, semester, alias, bio, avatar_url, github_username, github_id, github_token, role)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET full_name = $2,
    major = $3,
    semester = $4,
    alias = $5,
    bio = $6,
    avatar_url = $7,
    role = $8,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: MarkTokenInvalid :exec
UPDATE users
SET github_token = '',
    updated_at = NOW()
WHERE id = $1;
