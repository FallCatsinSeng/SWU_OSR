CREATE TABLE activity_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    showcase_repo_id UUID NOT NULL REFERENCES showcase_repos(id),
    event_type      VARCHAR(30) NOT NULL,
    summary         TEXT NOT NULL,
    metadata        JSONB NOT NULL DEFAULT '{}',
    github_event_id VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_created ON activity_logs(created_at DESC);
CREATE INDEX idx_activity_logs_user ON activity_logs(user_id, created_at DESC);
CREATE INDEX idx_activity_logs_event_type ON activity_logs(event_type);
