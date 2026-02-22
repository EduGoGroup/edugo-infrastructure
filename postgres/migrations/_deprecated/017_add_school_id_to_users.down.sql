-- Rollback: 017_add_school_id_to_users

DROP INDEX IF EXISTS idx_users_school_id;
ALTER TABLE users DROP COLUMN IF EXISTS school_id;
