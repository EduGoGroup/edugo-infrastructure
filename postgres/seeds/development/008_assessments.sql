-- =============================================================================
-- EduGo Development Seeds — 008_assessments.sql
-- =============================================================================
-- Crea 3 assessments de ejemplo con los nuevos campos boolean y date.
-- La relacion con materiales es N:N via assessment_materials.
--
-- NOTA sobre mongo_document_id:
--   VARCHAR(24) — simula el ObjectId de MongoDB que genera el worker.
--   Debe ser exactamente 24 caracteres hexadecimales.
--
-- NOTA sobre status (CHECK constraint):
--   'draft' | 'generated' | 'published' | 'archived' | 'closed'
--
-- Mapa de assessments:
--   ass001 → "Examen Fracciones"      — published, is_timed, shuffle, show_correct
--   ass002 → "Quiz Sistema Solar"     — published, sin timer, sin shuffle, sin show_correct
--   ass003 → "Evaluacion Historia"    — draft, is_timed, available_from futuro
--
-- Mapa de assessment_materials:
--   ass001 → mat001 (Fracciones) + mat002 (Sistema Solar)
--   ass002 → mat002 (Sistema Solar)
--   ass003 → (sin materiales, es draft)
-- =============================================================================

BEGIN;

-- =========================================================================
-- PARTE 1: Assessments
-- =========================================================================

INSERT INTO assessment.assessment (
    id,
    mongo_document_id,
    school_id,
    created_by_user_id,
    title,
    description,
    questions_count,
    pass_threshold,
    max_attempts,
    time_limit_minutes,
    is_timed,
    shuffle_questions,
    show_correct_answers,
    available_from,
    available_until,
    status
) VALUES

-- -------------------------------------------------------------------------
-- ass001: Examen Fracciones (published, timed, shuffle, show_correct)
-- -------------------------------------------------------------------------
(
    'aa200000-0000-0000-0000-000000000001',
    'aaaaaa0000000000000000a1',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000005',
    'Examen Fracciones',
    'Evaluacion sobre operaciones basicas con fracciones: suma, resta y equivalencias.',
    5,
    60,
    3,
    30,
    true,
    true,
    true,
    NOW() - INTERVAL '7 days',
    NOW() + INTERVAL '30 days',
    'published'
),

-- -------------------------------------------------------------------------
-- ass002: Quiz Sistema Solar (published, sin timer, sin shuffle)
-- -------------------------------------------------------------------------
(
    'aa200000-0000-0000-0000-000000000002',
    'aaaaaa0000000000000000a2',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000006',
    'Quiz Sistema Solar',
    'Quiz rapido sobre los planetas del sistema solar y sus caracteristicas.',
    5,
    60,
    2,
    25,
    false,
    false,
    false,
    NULL,
    NULL,
    'published'
),

-- -------------------------------------------------------------------------
-- ass003: Evaluacion Historia (draft, timed, available_from futuro)
-- -------------------------------------------------------------------------
(
    'aa200000-0000-0000-0000-000000000003',
    'aaaaaa0000000000000000a3',
    'b2000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000007',
    'Evaluacion Historia',
    'Evaluacion sobre los principales procesos historicos de America Latina.',
    0,
    70,
    NULL,
    45,
    true,
    false,
    true,
    NOW() + INTERVAL '7 days',
    NULL,
    'draft'
)

ON CONFLICT (mongo_document_id) DO UPDATE SET
    title                = EXCLUDED.title,
    description          = EXCLUDED.description,
    questions_count      = EXCLUDED.questions_count,
    pass_threshold       = EXCLUDED.pass_threshold,
    max_attempts         = EXCLUDED.max_attempts,
    time_limit_minutes   = EXCLUDED.time_limit_minutes,
    is_timed             = EXCLUDED.is_timed,
    shuffle_questions    = EXCLUDED.shuffle_questions,
    show_correct_answers = EXCLUDED.show_correct_answers,
    available_from       = EXCLUDED.available_from,
    available_until      = EXCLUDED.available_until,
    status               = EXCLUDED.status,
    updated_at           = now();

-- =========================================================================
-- PARTE 2: Assessment Materials (relacion N:N)
-- =========================================================================

INSERT INTO assessment.assessment_materials (
    id,
    assessment_id,
    material_id,
    sort_order
) VALUES

-- ass001 -> mat001 (Fracciones)
(
    'ab100000-0000-0000-0000-000000000001',
    'aa200000-0000-0000-0000-000000000001',
    'aa100000-0000-0000-0000-000000000001',
    0
),

-- ass001 -> mat002 (Sistema Solar)
(
    'ab100000-0000-0000-0000-000000000002',
    'aa200000-0000-0000-0000-000000000001',
    'aa100000-0000-0000-0000-000000000002',
    1
),

-- ass002 -> mat002 (Sistema Solar)
(
    'ab100000-0000-0000-0000-000000000003',
    'aa200000-0000-0000-0000-000000000002',
    'aa100000-0000-0000-0000-000000000002',
    0
)

ON CONFLICT (assessment_id, material_id) DO UPDATE SET
    sort_order = EXCLUDED.sort_order;

COMMIT;
