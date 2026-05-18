CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nim             VARCHAR(20) UNIQUE NOT NULL,
    full_name       VARCHAR(255) NOT NULL,
    major           VARCHAR(255) NOT NULL,
    semester        INT NOT NULL,
    alias           VARCHAR(50) UNIQUE NOT NULL,
    bio             TEXT DEFAULT '',
    avatar_url      VARCHAR(512) DEFAULT '',
    github_username VARCHAR(100) UNIQUE NOT NULL,
    github_id       BIGINT UNIQUE NOT NULL,
    github_token    TEXT NOT NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'student',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_users_alias ON users(alias) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_github_username ON users(github_username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
