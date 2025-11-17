-- Rollback: 011_extend_assessment_answer

BEGIN;

ALTER TABLE assessment_attempt_answer
    DROP COLUMN IF EXISTS time_spent_seconds;

COMMIT;
