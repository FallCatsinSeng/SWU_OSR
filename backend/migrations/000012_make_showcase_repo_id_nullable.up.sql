-- Allow activity_logs to record public activity from repos not in the showcase.
ALTER TABLE activity_logs ALTER COLUMN showcase_repo_id DROP NOT NULL;

-- Add repo_full_name column so we can display repo info even without a showcase entry.
ALTER TABLE activity_logs ADD COLUMN repo_full_name VARCHAR(255) NOT NULL DEFAULT '';

-- Backfill repo_full_name from existing showcase_repos references.
UPDATE activity_logs a
SET repo_full_name = s.repo_full_name
FROM showcase_repos s
WHERE a.showcase_repo_id = s.id AND a.repo_full_name = '';

-- Index for general feed queries that include non-showcase activity.
CREATE INDEX idx_activity_logs_repo_full_name ON activity_logs(repo_full_name) WHERE repo_full_name != '';
