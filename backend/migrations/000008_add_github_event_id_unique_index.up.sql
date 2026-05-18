CREATE UNIQUE INDEX idx_activity_logs_event_id ON activity_logs(github_event_id) WHERE github_event_id IS NOT NULL;
