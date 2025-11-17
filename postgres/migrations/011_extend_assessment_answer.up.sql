-- Migration: 011_extend_assessment_answer
-- Description: Extender assessment_attempt_answer con time_spent_seconds
-- Dependencies: 008_create_assessment_answers.up.sql
-- Date: 2025-11-17

BEGIN;

-- 1. Agregar campo time_spent_seconds
ALTER TABLE assessment_attempt_answer
    ADD COLUMN IF NOT EXISTS time_spent_seconds INTEGER CHECK (time_spent_seconds >= 0);

-- 2. Agregar comentario
COMMENT ON COLUMN assessment_attempt_answer.time_spent_seconds IS 'Tiempo que tomó responder esta pregunta en segundos';

-- 3. Crear alias columns para compatibilidad con isolated design
-- Nota: Mantenemos question_index (INTEGER actual) como fuente de verdad
-- vs question_id (VARCHAR isolated). Mapeo se hace en capa de aplicación.
-- Mantenemos student_answer (TEXT actual) como fuente de verdad
-- vs selected_answer_id (VARCHAR isolated). Mapeo se hace en capa de aplicación.

COMMENT ON COLUMN assessment_attempt_answer.question_index IS 'Índice de la pregunta (0-based). APIs mapean a question_id según necesidad.';
COMMENT ON COLUMN assessment_attempt_answer.student_answer IS 'Respuesta del estudiante (TEXT flexible: JSON, string, etc). APIs mapean a selected_answer_id según necesidad.';

COMMIT;
