-- Rollback: 009_extend_assessment_schema
-- WARNING: Esta operación elimina columnas y datos

BEGIN;

-- Drop trigger y función
DROP TRIGGER IF EXISTS trg_sync_questions_count ON assessment;
DROP FUNCTION IF EXISTS sync_questions_count();

-- Drop columnas agregadas
ALTER TABLE assessment
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS pass_threshold,
    DROP COLUMN IF EXISTS max_attempts,
    DROP COLUMN IF EXISTS time_limit_minutes,
    DROP COLUMN IF EXISTS total_questions;

-- Restaurar constraint original de status
ALTER TABLE assessment
    DROP CONSTRAINT IF EXISTS assessment_status_check,
    ADD CONSTRAINT assessment_status_check 
        CHECK (status IN ('generated', 'published', 'archived'));

COMMIT;
