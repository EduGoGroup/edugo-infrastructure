-- =============================================================================
-- EduGo Development Seeds v2 — 008_assessments.sql
-- =============================================================================
-- 8 assessments (5 published + 3 draft) across 3 schools.
-- ass001-006: ai_generated (mongo_document_id references MongoDB)
-- ass007-008: manual (questions stored in PG, no MongoDB)
--
-- Mapa:
--   ass001 → Examen Fracciones       → San Ignacio → Maria (U-05) → published, timed 30min, ai_generated
--   ass002 → Quiz Sistema Solar      → San Ignacio → Pedro (U-06) → published, no timer, ai_generated
--   ass003 → Ejercicio Color y Forma → CreArte     → Ana (U-07)   → published, no timer, ai_generated
--   ass004 → English Grammar Test    → Academia    → Maria (U-05) → published, timed 20min, ai_generated
--   ass005 → Evaluacion Historia     → San Ignacio → Pedro (U-06) → draft, ai_generated
--   ass006 → Proyecto Final Escultura→ CreArte     → Ana (U-07)   → draft, ai_generated
--   ass007 → Evaluacion Manual Mate  → San Ignacio → Maria (U-05) → published, manual, 4 preguntas PG
--   ass008 → Quiz Manual Ciencias    → San Ignacio → Maria (U-05) → draft, manual, 2 preguntas PG
--
-- Assessment-Material mapping:
--   ass001 → mat001 (Fracciones)
--   ass002 → mat002 (Sistema Solar)
--   ass003 → mat004 (Teoria del Color)
--   ass004 → mat005 (English Grammar)
--   ass005 → mat003 (Historia Chile)
--   ass006-008 → (ninguno)
-- =============================================================================

BEGIN;

-- =========================================================================
-- PARTE 1: Assessments
-- =========================================================================

INSERT INTO assessment.assessment (
    id, mongo_document_id, source_type, school_id, created_by_user_id,
    title, description, questions_count, pass_threshold, max_attempts,
    time_limit_minutes, is_timed, shuffle_questions, show_correct_answers,
    available_from, available_until, status
) VALUES

-- ass001: Examen Fracciones (published, timed 30min, 5 preguntas, pass 60%, ai_generated)
(
    'aa200000-0000-0000-0000-000000000001',
    'aaaaaa000000000000000001', 'ai_generated',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000005',
    'Examen Fracciones',
    'Evaluacion sobre operaciones basicas con fracciones: suma, resta y equivalencias.',
    5, 60, 3, 30, true, true, true,
    NOW() - INTERVAL '7 days',
    NOW() + INTERVAL '30 days',
    'published'
),

-- ass002: Quiz Sistema Solar (published, no timer, 4 preguntas, pass 50%, ai_generated)
(
    'aa200000-0000-0000-0000-000000000002',
    'aaaaaa000000000000000002', 'ai_generated',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000006',
    'Quiz Ciencias: Sistema Solar',
    'Quiz rapido sobre los planetas del sistema solar y sus caracteristicas.',
    4, 50, 2, NULL, false, false, false,
    NULL, NULL,
    'published'
),

-- ass003: Ejercicio Color y Forma (published, no timer, 3 preguntas, pass 70%, ai_generated)
(
    'aa200000-0000-0000-0000-000000000003',
    'aaaaaa000000000000000003', 'ai_generated',
    'b2000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000007',
    'Ejercicio Color y Forma',
    'Ejercicio practico sobre teoria del color y composicion visual.',
    3, 70, 2, NULL, false, false, true,
    NOW() - INTERVAL '5 days',
    NOW() + INTERVAL '60 days',
    'published'
),

-- ass004: English Grammar Test (published, timed 20min, 4 preguntas, pass 60%, ai_generated)
(
    'aa200000-0000-0000-0000-000000000004',
    'aaaaaa000000000000000004', 'ai_generated',
    'b3000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000005',
    'English Grammar Test',
    'Test on basic English grammar: articles, pronouns, and simple tenses.',
    4, 60, 2, 20, true, true, true,
    NOW() - INTERVAL '3 days',
    NOW() + INTERVAL '30 days',
    'published'
),

-- ass005: Evaluacion Historia Chile (draft, 3 preguntas, ai_generated)
(
    'aa200000-0000-0000-0000-000000000005',
    'aaaaaa000000000000000005', 'ai_generated',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000006',
    'Evaluacion Historia Chile',
    'Evaluacion sobre los principales procesos de la independencia de Chile.',
    3, 70, NULL, 45, true, false, true,
    NOW() + INTERVAL '7 days',
    NULL,
    'draft'
),

-- ass006: Proyecto Final Escultura (draft, 0 preguntas, ai_generated)
(
    'aa200000-0000-0000-0000-000000000006',
    'aaaaaa000000000000000006', 'ai_generated',
    'b2000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000007',
    'Proyecto Final Escultura',
    'Proyecto final del modulo de escultura: crear una pieza original.',
    0, 60, 1, NULL, false, false, false,
    NULL, NULL,
    'draft'
),

-- ass007: Evaluacion Manual Matematicas (published, manual, 4 preguntas PG, assigned)
(
    'aa200000-0000-0000-0000-000000000007',
    NULL, 'manual',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000005',
    'Evaluacion Manual: Operaciones Basicas',
    'Evaluacion manual creada por la profesora Maria. Incluye 4 tipos de pregunta: opcion multiple, verdadero/falso, respuesta corta y abierta.',
    4, 60, 2, 25, true, false, true,
    NOW() - INTERVAL '2 days',
    NOW() + INTERVAL '14 days',
    'published'
),

-- ass008: Quiz Manual Ciencias (draft, manual, 2 preguntas PG)
(
    'aa200000-0000-0000-0000-000000000008',
    NULL, 'manual',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000005',
    'Quiz Manual: Ciencias Naturales',
    'Quiz manual en borrador sobre fotosintesis y ecosistemas.',
    2, 50, 1, NULL, false, false, true,
    NULL, NULL,
    'draft'
)

ON CONFLICT (id) DO UPDATE SET
    title                = EXCLUDED.title,
    description          = EXCLUDED.description,
    source_type          = EXCLUDED.source_type,
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
    assessment_id, material_id, sort_order
) VALUES
-- ass001 → mat001 (Fracciones)
(
    'aa200000-0000-0000-0000-000000000001',
    'aa100000-0000-0000-0000-000000000001',
    0
),
-- ass002 → mat002 (Sistema Solar)
(
    'aa200000-0000-0000-0000-000000000002',
    'aa100000-0000-0000-0000-000000000002',
    0
),
-- ass003 → mat004 (Teoria del Color)
(
    'aa200000-0000-0000-0000-000000000003',
    'aa100000-0000-0000-0000-000000000004',
    0
),
-- ass004 → mat005 (English Grammar)
(
    'aa200000-0000-0000-0000-000000000004',
    'aa100000-0000-0000-0000-000000000005',
    0
),
-- ass005 → mat003 (Historia Chile)
(
    'aa200000-0000-0000-0000-000000000005',
    'aa100000-0000-0000-0000-000000000003',
    0
)
ON CONFLICT (assessment_id, material_id) DO UPDATE SET
    sort_order = EXCLUDED.sort_order;

-- =========================================================================
-- PARTE 3: Preguntas PG para assessments manuales (ass007, ass008)
-- =========================================================================

-- ass007: 4 preguntas (MC, T/F, short_answer, open_ended)
INSERT INTO assessment.questions (
    id, assessment_id, question_text, question_type, correct_answer,
    explanation, points, difficulty, sort_order
) VALUES
-- Q1: multiple_choice
(
    'aq000000-0000-0000-0000-000000000001',
    'aa200000-0000-0000-0000-000000000007',
    'Cuanto es 15 + 27?',
    'multiple_choice', '42',
    'Se suman las unidades (5+7=12, llevamos 1) y las decenas (1+2+1=4). Resultado: 42.',
    2.0, 'easy', 0
),
-- Q2: true_false
(
    'aq000000-0000-0000-0000-000000000002',
    'aa200000-0000-0000-0000-000000000007',
    'El resultado de multiplicar cualquier numero por cero es cero.',
    'true_false', 'Verdadero',
    'La propiedad absorbente de la multiplicacion establece que a x 0 = 0 para todo numero a.',
    1.0, 'easy', 1
),
-- Q3: short_answer
(
    'aq000000-0000-0000-0000-000000000003',
    'aa200000-0000-0000-0000-000000000007',
    'Como se llama el resultado de una resta?',
    'short_answer', 'diferencia',
    'El resultado de una resta se denomina diferencia.',
    1.5, 'medium', 2
),
-- Q4: open_ended
(
    'aq000000-0000-0000-0000-000000000004',
    'aa200000-0000-0000-0000-000000000007',
    'Explica con un ejemplo de la vida cotidiana donde usarias la multiplicacion.',
    'open_ended', NULL,
    'Respuesta abierta. Se evalua la capacidad de relacionar operaciones matematicas con situaciones reales.',
    5.0, 'hard', 3
)
ON CONFLICT (id) DO NOTHING;

-- ass008: 2 preguntas (MC, short_answer) — draft
INSERT INTO assessment.questions (
    id, assessment_id, question_text, question_type, correct_answer,
    explanation, points, difficulty, sort_order
) VALUES
(
    'aq000000-0000-0000-0000-000000000005',
    'aa200000-0000-0000-0000-000000000008',
    'Que gas absorben las plantas durante la fotosintesis?',
    'multiple_choice', 'Dioxido de carbono',
    'Las plantas absorben CO2 y liberan O2 durante la fotosintesis.',
    2.0, 'easy', 0
),
(
    'aq000000-0000-0000-0000-000000000006',
    'aa200000-0000-0000-0000-000000000008',
    'Como se llama la capa de la atmosfera donde vivimos?',
    'short_answer', 'troposfera',
    'La troposfera es la capa mas baja de la atmosfera terrestre.',
    2.0, 'medium', 1
)
ON CONFLICT (id) DO NOTHING;

-- =========================================================================
-- PARTE 4: Opciones para preguntas multiple_choice
-- =========================================================================

-- Opciones Q1 (15+27): 32, 42, 52
INSERT INTO assessment.question_options (id, question_id, option_text, sort_order) VALUES
('ao000000-0000-0000-0000-000000000001', 'aq000000-0000-0000-0000-000000000001', '32', 0),
('ao000000-0000-0000-0000-000000000002', 'aq000000-0000-0000-0000-000000000001', '42', 1),
('ao000000-0000-0000-0000-000000000003', 'aq000000-0000-0000-0000-000000000001', '52', 2),
('ao000000-0000-0000-0000-000000000004', 'aq000000-0000-0000-0000-000000000001', '41', 3)
ON CONFLICT (id) DO NOTHING;

-- Opciones Q5 (gas fotosintesis): Oxigeno, CO2, Nitrogeno
INSERT INTO assessment.question_options (id, question_id, option_text, sort_order) VALUES
('ao000000-0000-0000-0000-000000000005', 'aq000000-0000-0000-0000-000000000005', 'Oxigeno', 0),
('ao000000-0000-0000-0000-000000000006', 'aq000000-0000-0000-0000-000000000005', 'Dioxido de carbono', 1),
('ao000000-0000-0000-0000-000000000007', 'aq000000-0000-0000-0000-000000000005', 'Nitrogeno', 2)
ON CONFLICT (id) DO NOTHING;

-- =========================================================================
-- PARTE 5: Asignaciones (ass007 published → asignada a estudiantes)
-- =========================================================================

INSERT INTO assessment.assessment_assignments (
    id, assessment_id, student_id, assigned_by, assigned_at
) VALUES
-- Carlos Mendoza (est.carlos) asignado por Maria Martinez
(
    'ab000000-0000-0000-0000-000000000001',
    'aa200000-0000-0000-0000-000000000007',
    '00000000-0000-0000-0000-000000000008',
    '00000000-0000-0000-0000-000000000005',
    NOW() - INTERVAL '1 day'
),
-- Sofia Lopez (est.sofia) asignada por Maria Martinez
(
    'ab000000-0000-0000-0000-000000000002',
    'aa200000-0000-0000-0000-000000000007',
    '00000000-0000-0000-0000-000000000009',
    '00000000-0000-0000-0000-000000000005',
    NOW() - INTERVAL '1 day'
),
-- Diego Ramirez (est.diego) asignado por Maria Martinez
(
    'ab000000-0000-0000-0000-000000000003',
    'aa200000-0000-0000-0000-000000000007',
    '00000000-0000-0000-0000-000000000010',
    '00000000-0000-0000-0000-000000000005',
    NOW() - INTERVAL '1 day'
)
ON CONFLICT (id) DO NOTHING;

COMMIT;
