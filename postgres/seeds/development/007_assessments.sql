-- =============================================================================
-- EduGo Development Seeds — 007_assessments.sql
-- =============================================================================
-- Crea 2 assessments generados a partir de los materiales procesados.
-- Solo mat001 y mat002 tienen assessment (mat003 y mat004 no aplican aún).
--
-- NOTA sobre el trigger trg_sync_questions_count:
--   El trigger sincroniza questions_count con total_questions en INSERT/UPDATE.
--   Se insertan ambos con el mismo valor para coherencia.
--
-- NOTA sobre mongo_document_id:
--   VARCHAR(24) — simula el ObjectId de MongoDB que genera el worker.
--   Debe ser exactamente 24 caracteres hexadecimales.
--
-- NOTA sobre status (CHECK constraint):
--   'draft' | 'generated' | 'published' | 'archived' | 'closed'
--   Se usa 'published' para que los estudiantes puedan intentarlo.
--
-- Mapa de assessments:
--   ass001 → "Evaluación: Fracciones"    — mat001, 5 preguntas, umbral 60%
--   ass002 → "Evaluación: Sistema Solar" — mat002, 5 preguntas, umbral 60%
-- =============================================================================

BEGIN;

INSERT INTO public.assessment (
    id,
    material_id,
    mongo_document_id,
    title,
    questions_count,
    total_questions,
    pass_threshold,
    max_attempts,
    time_limit_minutes,
    status
) VALUES

-- -------------------------------------------------------------------------
-- Assessment de Fracciones (material mat001)
-- -------------------------------------------------------------------------
(
    'aa200000-0000-0000-0000-000000000001',
    'aa100000-0000-0000-0000-000000000001',   -- Introducción a las Fracciones
    'aaaaaa0000000000000000a1',               -- 24 chars hex (ObjectId simulado para mat001)
    'Evaluación: Fracciones',
    5,
    5,
    60,
    3,
    30,
    'published'
),

-- -------------------------------------------------------------------------
-- Assessment de Sistema Solar (material mat002)
-- -------------------------------------------------------------------------
(
    'aa200000-0000-0000-0000-000000000002',
    'aa100000-0000-0000-0000-000000000002',   -- El Sistema Solar
    'aaaaaa0000000000000000a2',               -- 24 chars hex (ObjectId simulado para mat002)
    'Evaluación: Sistema Solar',
    5,
    5,
    60,
    2,
    25,
    'published'
)

ON CONFLICT (mongo_document_id) DO UPDATE SET
    title               = EXCLUDED.title,
    questions_count     = EXCLUDED.questions_count,
    total_questions     = EXCLUDED.total_questions,
    pass_threshold      = EXCLUDED.pass_threshold,
    max_attempts        = EXCLUDED.max_attempts,
    time_limit_minutes  = EXCLUDED.time_limit_minutes,
    status              = EXCLUDED.status,
    updated_at          = now();

COMMIT;
