-- Migration: 009_extend_assessment_schema
-- Description: Extender schema de assessment con campos del diseño isolated
-- Dependencies: 006_create_assessments.up.sql
-- Date: 2025-11-17

BEGIN;

-- 1. Agregar campos faltantes a assessment (sin CHECK constraints inline)
ALTER TABLE assessment
    ADD COLUMN IF NOT EXISTS title VARCHAR(255),
    ADD COLUMN IF NOT EXISTS pass_threshold INTEGER DEFAULT 70,
    ADD COLUMN IF NOT EXISTS max_attempts INTEGER,
    ADD COLUMN IF NOT EXISTS time_limit_minutes INTEGER;

-- 1.1 Agregar CHECK constraint de forma idempotente
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'assessment_pass_threshold_check'
          AND conrelid = 'assessment'::regclass
    ) THEN
        ALTER TABLE assessment
            ADD CONSTRAINT assessment_pass_threshold_check
            CHECK (pass_threshold >= 0 AND pass_threshold <= 100);
    END IF;
END
$$;

-- 2. Agregar total_questions (sincronizado con questions_count)
ALTER TABLE assessment
    ADD COLUMN IF NOT EXISTS total_questions INTEGER;

-- Sincronizar total_questions = questions_count para datos existentes
UPDATE assessment SET total_questions = questions_count WHERE total_questions IS NULL;

-- 3. Agregar comentarios
COMMENT ON COLUMN assessment.title IS 'Título del assessment (opcional si se usa metadata de MongoDB)';
COMMENT ON COLUMN assessment.pass_threshold IS 'Porcentaje mínimo para aprobar (0-100)';
COMMENT ON COLUMN assessment.max_attempts IS 'Máximo de intentos permitidos (NULL = ilimitado)';
COMMENT ON COLUMN assessment.time_limit_minutes IS 'Límite de tiempo en minutos (NULL = sin límite)';
COMMENT ON COLUMN assessment.total_questions IS 'Total de preguntas (sincronizado con questions_count)';

-- 4. Actualizar status values (agregar 'draft' y 'closed')
-- Nota: No eliminamos valores existentes para retrocompatibilidad
ALTER TABLE assessment
    DROP CONSTRAINT IF EXISTS assessment_status_check,
    ADD CONSTRAINT assessment_status_check 
        CHECK (status IN ('draft', 'generated', 'published', 'archived', 'closed'));

-- 5. Trigger para mantener sincronizado questions_count y total_questions (transición gradual)
CREATE OR REPLACE FUNCTION sync_questions_count()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.total_questions IS NOT NULL THEN
        NEW.questions_count := NEW.total_questions;
    ELSIF NEW.questions_count IS NOT NULL THEN
        NEW.total_questions := NEW.questions_count;
    ELSE
        -- Si ambos son NULL, inicializar en 0
        NEW.total_questions := 0;
        NEW.questions_count := 0;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_sync_questions_count
    BEFORE INSERT OR UPDATE ON assessment
    FOR EACH ROW
    EXECUTE FUNCTION sync_questions_count();

COMMENT ON TRIGGER trg_sync_questions_count ON assessment IS 'Mantiene sincronizado questions_count y total_questions durante transición';

COMMIT;
