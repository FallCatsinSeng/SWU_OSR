CREATE TABLE showcase_repos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    github_repo_id  BIGINT NOT NULL,
    repo_name       VARCHAR(255) NOT NULL,
    repo_full_name  VARCHAR(512) NOT NULL,
    description     TEXT DEFAULT '',
    language        VARCHAR(50) DEFAULT '',
    html_url        VARCHAR(512) NOT NULL,
    academic_tag    VARCHAR(50) NOT NULL,
    webhook_id      BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE(user_id, github_repo_id)
);

CREATE INDEX idx_showcase_repos_user ON showcase_repos(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_showcase_repos_tag ON showcase_repos(academic_tag) WHERE deleted_at IS NULL;
