-- Migration: 010_extend_assessment_attempt
-- Description: Extender assessment_attempt con time_spent_seconds e idempotency_key
-- Dependencies: 007_create_assessment_attempts.up.sql
-- Date: 2025-11-17

BEGIN;

-- 1. Agregar campos faltantes
ALTER TABLE assessment_attempt
    ADD COLUMN IF NOT EXISTS time_spent_seconds INTEGER CHECK (time_spent_seconds > 0 AND time_spent_seconds <= 7200),
    ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(64);

-- 2. Agregar UNIQUE constraint a idempotency_key (separado para evitar conflictos)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'unique_idempotency_key'
    ) THEN
        ALTER TABLE assessment_attempt
            ADD CONSTRAINT unique_idempotency_key UNIQUE (idempotency_key);
    END IF;
END $$;

-- 3. Agregar comentarios
COMMENT ON COLUMN assessment_attempt.time_spent_seconds IS 'Tiempo total del intento en segundos (max 2 horas)';
COMMENT ON COLUMN assessment_attempt.idempotency_key IS 'Clave para prevenir intentos duplicados';

-- 4. Crear índice parcial para idempotency_key
CREATE INDEX IF NOT EXISTS idx_attempt_idempotency_key 
    ON assessment_attempt(idempotency_key) 
    WHERE idempotency_key IS NOT NULL;

-- 5. Agregar CHECK constraints
ALTER TABLE assessment_attempt
    ADD CONSTRAINT IF NOT EXISTS check_attempt_time_logical 
        CHECK (completed_at IS NULL OR completed_at > started_at);

COMMENT ON CONSTRAINT check_attempt_time_logical ON assessment_attempt IS 'Validar que completed_at > started_at';

-- Nota: No agregamos constraint de time_spent = (completed - started) 
-- porque datos existentes pueden no cumplirlo. Se validará en aplicación.

COMMIT;
