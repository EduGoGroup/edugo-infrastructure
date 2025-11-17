-- Rollback: 010_extend_assessment_attempt

BEGIN;

-- Drop índice
DROP INDEX IF EXISTS idx_attempt_idempotency_key;

-- Drop constraints
ALTER TABLE assessment_attempt
    DROP CONSTRAINT IF EXISTS check_attempt_time_logical,
    DROP CONSTRAINT IF EXISTS unique_idempotency_key;

-- Drop columnas
ALTER TABLE assessment_attempt
    DROP COLUMN IF EXISTS time_spent_seconds,
    DROP COLUMN IF EXISTS idempotency_key;

COMMIT;
