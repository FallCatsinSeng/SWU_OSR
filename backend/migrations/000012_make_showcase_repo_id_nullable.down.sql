-- Remove non-showcase activity entries before re-adding NOT NULL constraint.
DELETE FROM activity_logs WHERE showcase_repo_id IS NULL;

DROP INDEX IF EXISTS idx_activity_logs_repo_full_name;
ALTER TABLE activity_logs DROP COLUMN IF EXISTS repo_full_name;
ALTER TABLE activity_logs ALTER COLUMN showcase_repo_id SET NOT NULL;
