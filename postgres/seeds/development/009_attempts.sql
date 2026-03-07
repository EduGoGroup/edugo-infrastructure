-- =============================================================================
-- EduGo Development Seeds v2 — 009_attempts.sql
-- =============================================================================
-- 7 intentos de evaluacion + respuestas detalladas del primer intento de Carlos.
--
-- Mapa:
--   att001 → Carlos (U-08)    → Fracciones (ass001) → score 80/100  → intento 1
--   att002 → Carlos (U-08)    → Fracciones (ass001) → score 92/100  → intento 2
--   att003 → Sofia (U-09)     → Fracciones (ass001) → score 68/100  → intento 1
--   att004 → Diego (U-10)     → Ciencias (ass002)   → score 75/80   → intento 1
--   att005 → Carlos (U-08)    → Color y Forma (ass003) → score 60/100 → intento 1
--   att006 → Mateo (U-12)     → Color y Forma (ass003) → score 90/100 → intento 1
--   att007 → Valentina (U-11) → English Grammar (ass004) → score 85/100 → intento 1
-- =============================================================================

BEGIN;

INSERT INTO assessment.assessment_attempt (
    id, assessment_id, student_id,
    started_at, completed_at,
    score, max_score, percentage,
    status, time_spent_seconds, idempotency_key
) VALUES

-- att001: Carlos → Fracciones intento 1
(
    'aa300000-0000-0000-0000-000000000001',
    'aa200000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000008',
    NOW() - INTERVAL '3 days',
    NOW() - INTERVAL '3 days' + INTERVAL '25 minutes',
    80.00, 100.00, 80.00,
    'completed', 1520,
    'idem_att001_carlos_ass001_v2'
),

-- att002: Carlos → Fracciones intento 2
(
    'aa300000-0000-0000-0000-000000000002',
    'aa200000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000008',
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '2 days' + INTERVAL '20 minutes',
    92.00, 100.00, 92.00,
    'completed', 1200,
    'idem_att002_carlos_ass001_v2'
),

-- att003: Sofia → Fracciones intento 1
(
    'aa300000-0000-0000-0000-000000000003',
    'aa200000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000009',
    NOW() - INTERVAL '2 days' - INTERVAL '4 hours',
    NOW() - INTERVAL '2 days' - INTERVAL '3 hours 30 minutes',
    68.00, 100.00, 68.00,
    'completed', 1800,
    'idem_att003_sofia_ass001_v2'
),

-- att004: Diego → Ciencias intento 1
(
    'aa300000-0000-0000-0000-000000000004',
    'aa200000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000010',
    NOW() - INTERVAL '1 day' - INTERVAL '6 hours',
    NOW() - INTERVAL '1 day' - INTERVAL '5 hours 45 minutes',
    75.00, 80.00, 93.75,
    'completed', 900,
    'idem_att004_diego_ass002_v2'
),

-- att005: Carlos → Color y Forma intento 1 (en CreArte)
(
    'aa300000-0000-0000-0000-000000000005',
    'aa200000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000008',
    NOW() - INTERVAL '1 day' - INTERVAL '2 hours',
    NOW() - INTERVAL '1 day' - INTERVAL '1 hour 40 minutes',
    60.00, 100.00, 60.00,
    'completed', 1200,
    'idem_att005_carlos_ass003_v2'
),

-- att006: Mateo → Color y Forma intento 1 (en CreArte)
(
    'aa300000-0000-0000-0000-000000000006',
    'aa200000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000012',
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '23 hours',
    90.00, 100.00, 90.00,
    'completed', 600,
    'idem_att006_mateo_ass003_v2'
),

-- att007: Valentina → English Grammar intento 1 (en Academia)
(
    'aa300000-0000-0000-0000-000000000007',
    'aa200000-0000-0000-0000-000000000004',
    '00000000-0000-0000-0000-000000000011',
    NOW() - INTERVAL '12 hours',
    NOW() - INTERVAL '11 hours 45 minutes',
    85.00, 100.00, 85.00,
    'completed', 1100,
    'idem_att007_valentina_ass004_v2'
)

ON CONFLICT (idempotency_key) DO UPDATE SET
    score              = EXCLUDED.score,
    max_score          = EXCLUDED.max_score,
    percentage         = EXCLUDED.percentage,
    status             = EXCLUDED.status,
    completed_at       = EXCLUDED.completed_at,
    time_spent_seconds = EXCLUDED.time_spent_seconds,
    updated_at         = now();

-- =========================================================================
-- Respuestas del intento 1 de Carlos (att001 → ass001 Fracciones)
-- 5 preguntas x 20 pts = 100 pts max
-- Correctas: q0, q1, q3, q4 (4/5) = 80 pts
-- =========================================================================

INSERT INTO assessment.assessment_attempt_answer (
    id, attempt_id, question_index,
    student_answer, is_correct, points_earned, max_points,
    time_spent_seconds, answered_at
) VALUES
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001', 0,
    '1/2', true, 20.00, 20.00, 280,
    NOW() - INTERVAL '3 days' + INTERVAL '5 minutes'
),
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001', 1,
    '3/4', true, 20.00, 20.00, 310,
    NOW() - INTERVAL '3 days' + INTERVAL '10 minutes'
),
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001', 2,
    '2/6', false, 0.00, 20.00, 420,
    NOW() - INTERVAL '3 days' + INTERVAL '17 minutes'
),
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001', 3,
    '2/5', true, 20.00, 20.00, 265,
    NOW() - INTERVAL '3 days' + INTERVAL '21 minutes'
),
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001', 4,
    '3/8', true, 20.00, 20.00, 245,
    NOW() - INTERVAL '3 days' + INTERVAL '25 minutes'
)
ON CONFLICT (attempt_id, question_index) DO UPDATE SET
    student_answer     = EXCLUDED.student_answer,
    is_correct         = EXCLUDED.is_correct,
    points_earned      = EXCLUDED.points_earned,
    max_points         = EXCLUDED.max_points,
    time_spent_seconds = EXCLUDED.time_spent_seconds,
    updated_at         = now();

-- =========================================================================
-- Attempt in_progress para testing del flujo de toma de evaluación
-- Carlos Mendoza (U-08) intentando "Examen Fracciones" (ass001)
-- =========================================================================

INSERT INTO assessment.assessment_attempt (
    id, assessment_id, student_id,
    started_at, completed_at,
    score, max_score, percentage,
    status, time_spent_seconds, idempotency_key
) VALUES
(
    'aa300000-0000-0000-0000-000000000010',
    'aa200000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000008',
    NOW() - INTERVAL '5 minutes',
    NULL,
    NULL, NULL, NULL,
    'in_progress', NULL,
    'idem_att010_carlos_ass001_inprogress'
)
ON CONFLICT (idempotency_key) DO UPDATE SET
    status             = EXCLUDED.status,
    completed_at       = EXCLUDED.completed_at,
    score              = EXCLUDED.score,
    updated_at         = now();

-- Respuestas parciales del attempt in_progress (2 de 5 preguntas)

INSERT INTO assessment.assessment_attempt_answer (
    id, attempt_id, question_index,
    student_answer, is_correct, points_earned, max_points,
    time_spent_seconds, answered_at
) VALUES
(
    'an000000-0000-0000-0000-000000000020',
    'aa300000-0000-0000-0000-000000000010', 0,
    'A', NULL, NULL, NULL, 5,
    NOW() - INTERVAL '4 minutes'
),
(
    'an000000-0000-0000-0000-000000000021',
    'aa300000-0000-0000-0000-000000000010', 1,
    'B', NULL, NULL, NULL, 8,
    NOW() - INTERVAL '3 minutes'
)
ON CONFLICT (attempt_id, question_index) DO UPDATE SET
    student_answer     = EXCLUDED.student_answer,
    is_correct         = EXCLUDED.is_correct,
    points_earned      = EXCLUDED.points_earned,
    max_points         = EXCLUDED.max_points,
    time_spent_seconds = EXCLUDED.time_spent_seconds,
    updated_at         = now();

COMMIT;
