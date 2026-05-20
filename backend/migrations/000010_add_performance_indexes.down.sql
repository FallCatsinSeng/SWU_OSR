-- Rollback performance indexes
DROP INDEX IF EXISTS idx_activity_logs_repo_created;
DROP INDEX IF EXISTS idx_activity_logs_user_event_created;
DROP INDEX IF EXISTS idx_threads_author_created;
DROP INDEX IF EXISTS idx_comments_author_created;
DROP INDEX IF EXISTS idx_showcase_repos_user_active;
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_activity_logs_showcase_repo_id;
