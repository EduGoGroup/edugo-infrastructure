-- =============================================================================
-- EduGo Development Seeds v2 — 008_assessments.sql
-- =============================================================================
-- 6 assessments (4 published + 2 draft) across 3 schools.
-- mongo_document_id: 24-char hex referencing MongoDB documents.
--
-- Mapa:
--   ass001 → Examen Fracciones       → San Ignacio → Maria (U-05) → published, timed 30min
--   ass002 → Quiz Sistema Solar      → San Ignacio → Pedro (U-06) → published, no timer
--   ass003 → Ejercicio Color y Forma → CreArte     → Ana (U-07)   → published, no timer
--   ass004 → English Grammar Test    → Academia    → Maria (U-05) → published, timed 20min
--   ass005 → Evaluacion Historia     → San Ignacio → Pedro (U-06) → draft
--   ass006 → Proyecto Final Escultura→ CreArte     → Ana (U-07)   → draft
--
-- Assessment-Material mapping:
--   ass001 → mat001 (Fracciones)
--   ass002 → mat002 (Sistema Solar)
--   ass003 → mat004 (Teoria del Color)
--   ass004 → mat005 (English Grammar)
--   ass005 → mat003 (Historia Chile)
--   ass006 → (ninguno)
-- =============================================================================

BEGIN;

-- =========================================================================
-- PARTE 1: Assessments
-- =========================================================================

INSERT INTO assessment.assessment (
    id, mongo_document_id, school_id, created_by_user_id,
    title, description, questions_count, pass_threshold, max_attempts,
    time_limit_minutes, is_timed, shuffle_questions, show_correct_answers,
    available_from, available_until, status
) VALUES

-- ass001: Examen Fracciones (published, timed 30min, 5 preguntas, pass 60%)
(
    'aa200000-0000-0000-0000-000000000001',
    'aaaaaa000000000000000001',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000005',
    'Examen Fracciones',
    'Evaluacion sobre operaciones basicas con fracciones: suma, resta y equivalencias.',
    5, 60, 3, 30, true, true, true,
    NOW() - INTERVAL '7 days',
    NOW() + INTERVAL '30 days',
    'published'
),

-- ass002: Quiz Sistema Solar (published, no timer, 4 preguntas, pass 50%)
(
    'aa200000-0000-0000-0000-000000000002',
    'aaaaaa000000000000000002',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000006',
    'Quiz Ciencias: Sistema Solar',
    'Quiz rapido sobre los planetas del sistema solar y sus caracteristicas.',
    4, 50, 2, NULL, false, false, false,
    NULL, NULL,
    'published'
),

-- ass003: Ejercicio Color y Forma (published, no timer, 3 preguntas, pass 70%)
(
    'aa200000-0000-0000-0000-000000000003',
    'aaaaaa000000000000000003',
    'b2000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000007',
    'Ejercicio Color y Forma',
    'Ejercicio practico sobre teoria del color y composicion visual.',
    3, 70, 2, NULL, false, false, true,
    NOW() - INTERVAL '5 days',
    NOW() + INTERVAL '60 days',
    'published'
),

-- ass004: English Grammar Test (published, timed 20min, 4 preguntas, pass 60%)
(
    'aa200000-0000-0000-0000-000000000004',
    'aaaaaa000000000000000004',
    'b3000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000005',
    'English Grammar Test',
    'Test on basic English grammar: articles, pronouns, and simple tenses.',
    4, 60, 2, 20, true, true, true,
    NOW() - INTERVAL '3 days',
    NOW() + INTERVAL '30 days',
    'published'
),

-- ass005: Evaluacion Historia Chile (draft, 3 preguntas)
(
    'aa200000-0000-0000-0000-000000000005',
    'aaaaaa000000000000000005',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000006',
    'Evaluacion Historia Chile',
    'Evaluacion sobre los principales procesos de la independencia de Chile.',
    3, 70, NULL, 45, true, false, true,
    NOW() + INTERVAL '7 days',
    NULL,
    'draft'
),

-- ass006: Proyecto Final Escultura (draft, 0 preguntas)
(
    'aa200000-0000-0000-0000-000000000006',
    'aaaaaa000000000000000006',
    'b2000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000007',
    'Proyecto Final Escultura',
    'Proyecto final del modulo de escultura: crear una pieza original.',
    0, 60, 1, NULL, false, false, false,
    NULL, NULL,
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
    id, assessment_id, material_id, sort_order
) VALUES
-- ass001 → mat001 (Fracciones)
(
    'ab100000-0000-0000-0000-000000000001',
    'aa200000-0000-0000-0000-000000000001',
    'aa100000-0000-0000-0000-000000000001',
    0
),
-- ass002 → mat002 (Sistema Solar)
(
    'ab100000-0000-0000-0000-000000000002',
    'aa200000-0000-0000-0000-000000000002',
    'aa100000-0000-0000-0000-000000000002',
    0
),
-- ass003 → mat004 (Teoria del Color)
(
    'ab100000-0000-0000-0000-000000000003',
    'aa200000-0000-0000-0000-000000000003',
    'aa100000-0000-0000-0000-000000000004',
    0
),
-- ass004 → mat005 (English Grammar)
(
    'ab100000-0000-0000-0000-000000000004',
    'aa200000-0000-0000-0000-000000000004',
    'aa100000-0000-0000-0000-000000000005',
    0
),
-- ass005 → mat003 (Historia Chile)
(
    'ab100000-0000-0000-0000-000000000005',
    'aa200000-0000-0000-0000-000000000005',
    'aa100000-0000-0000-0000-000000000003',
    0
)
ON CONFLICT (assessment_id, material_id) DO UPDATE SET
    sort_order = EXCLUDED.sort_order;

COMMIT;
