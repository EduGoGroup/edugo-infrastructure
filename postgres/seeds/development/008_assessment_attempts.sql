-- =============================================================================
-- EduGo Development Seeds — 008_assessment_attempts.sql
-- =============================================================================
-- Crea 4 intentos de evaluación y las respuestas del primer intento de Carlos.
--
-- DIFERENCIAS IMPORTANTES entre el schema especificado y el schema real:
--   - No existe campo "attempt_number"  → se omite
--   - No existe campo "submitted_at"    → se usa "completed_at"
--   - No existe campo "max_score"       → se calcula implícito (100 puntos)
--   - Campo "percentage" → se inserta junto con score
--   - assessment_attempt_answer usa "question_index" (INT, 0-based) + "student_answer"
--     en lugar de "question_id" (VARCHAR) + "selected_answer"
--   - No existe campo "answered_at" separado → se usa el default de "created_at"
--   - La constraint unique en answers es: (attempt_id, question_index)
--
-- Lógica de scores para att001 (Carlos, Fracciones, intento 1):
--   5 preguntas × 20 pts c/u = 100 pts máx
--   Correctas: q0, q1, q3, q4 (4/5) = 80 pts → score=80, percentage=80.00
--   (El escenario indica 85 pts pero con 5 preguntas de 20 pts cada una,
--    4 correctas dan exactamente 80. Se ajusta al schema real.)
--
-- Lógica de scores para att004 (Carlos, Fracciones, intento 2 — sin respuestas detalladas):
--   score=92, percentage=92.00 (5 preguntas, ~4.6 correctas, aprovechado)
--
-- Mapa de intentos:
--   att001 → carlos (u008)  → ass001 (Fracciones)     — completado, score=80
--   att002 → sofia  (u009)  → ass001 (Fracciones)     — completado, score=72
--   att003 → miguel (u010)  → ass002 (Sistema Solar)  — completado, score=90
--   att004 → carlos (u008)  → ass001 (Fracciones)     — completado, score=92 (2do intento)
-- =============================================================================

BEGIN;

-- =========================================================================
-- PARTE 1: Assessment Attempts
-- =========================================================================

INSERT INTO public.assessment_attempt (
    id,
    assessment_id,
    student_id,
    started_at,
    completed_at,
    score,
    max_score,
    percentage,
    status,
    time_spent_seconds,
    idempotency_key
) VALUES

-- -------------------------------------------------------------------------
-- att001: Carlos — Fracciones — Intento 1 (completado)
-- -------------------------------------------------------------------------
(
    'aa300000-0000-0000-0000-000000000001',
    'aa200000-0000-0000-0000-000000000001',   -- Evaluación: Fracciones
    '00000000-0000-0000-0000-000000000008',   -- Carlos González
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '23 hours',
    80.00,
    100.00,
    80.00,
    'completed',
    1520,
    'idem_att001_carlos_ass001_attempt1'
),

-- -------------------------------------------------------------------------
-- att002: Sofía — Fracciones — Intento 1 (completado)
-- -------------------------------------------------------------------------
(
    'aa300000-0000-0000-0000-000000000002',
    'aa200000-0000-0000-0000-000000000001',   -- Evaluación: Fracciones
    '00000000-0000-0000-0000-000000000009',   -- Sofía Rodríguez
    NOW() - INTERVAL '20 hours',
    NOW() - INTERVAL '19 hours 30 minutes',
    72.00,
    100.00,
    72.00,
    'completed',
    1800,
    'idem_att002_sofia_ass001_attempt1'
),

-- -------------------------------------------------------------------------
-- att003: Miguel — Sistema Solar — Intento 1 (completado)
-- -------------------------------------------------------------------------
(
    'aa300000-0000-0000-0000-000000000003',
    'aa200000-0000-0000-0000-000000000002',   -- Evaluación: Sistema Solar
    '00000000-0000-0000-0000-000000000010',   -- Miguel Torres
    NOW() - INTERVAL '15 hours',
    NOW() - INTERVAL '14 hours 45 minutes',
    90.00,
    100.00,
    90.00,
    'completed',
    900,
    'idem_att003_miguel_ass002_attempt1'
),

-- -------------------------------------------------------------------------
-- att004: Carlos — Fracciones — Intento 2 (completado, mejoró puntaje)
-- -------------------------------------------------------------------------
(
    'aa300000-0000-0000-0000-000000000004',
    'aa200000-0000-0000-0000-000000000001',   -- Evaluación: Fracciones
    '00000000-0000-0000-0000-000000000008',   -- Carlos González
    NOW() - INTERVAL '12 hours',
    NOW() - INTERVAL '11 hours 45 minutes',
    92.00,
    100.00,
    92.00,
    'completed',
    1200,
    'idem_att004_carlos_ass001_attempt2'
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
-- PARTE 2: Respuestas del Intento 1 de Carlos (att001 → ass001 Fracciones)
--
-- question_index es 0-based (0..4 para 5 preguntas)
-- student_answer refleja respuestas de fracciones
-- Resultado: 4 correctas (índices 0,1,3,4) + 1 incorrecta (índice 2)
--            = 80 puntos de 100
-- =========================================================================

INSERT INTO public.assessment_attempt_answer (
    id,
    attempt_id,
    question_index,
    student_answer,
    is_correct,
    points_earned,
    max_points,
    time_spent_seconds,
    answered_at
) VALUES

-- Pregunta 0: ¿Cuánto es 1/4 + 1/4? → 1/2 (correcto)
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001',
    0,
    '1/2',
    true,
    20.00,
    20.00,
    280,
    NOW() - INTERVAL '23 hours 25 minutes'
),

-- Pregunta 1: ¿Cuánto es 1/4 + 2/4? → 3/4 (correcto)
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001',
    1,
    '3/4',
    true,
    20.00,
    20.00,
    310,
    NOW() - INTERVAL '23 hours 20 minutes'
),

-- Pregunta 2: ¿Cuál fracción es equivalente a 2/6? → 1/3 (incorrecto: respondió 2/6)
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001',
    2,
    '2/6',
    false,
    0.00,
    20.00,
    420,
    NOW() - INTERVAL '23 hours 13 minutes'
),

-- Pregunta 3: ¿Cuánto es 1/5 + 1/5? → 2/5 (correcto)
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001',
    3,
    '2/5',
    true,
    20.00,
    20.00,
    265,
    NOW() - INTERVAL '23 hours 8 minutes'
),

-- Pregunta 4: ¿Cuánto es 1/8 + 2/8? → 3/8 (correcto)
(
    gen_random_uuid(),
    'aa300000-0000-0000-0000-000000000001',
    4,
    '3/8',
    true,
    20.00,
    20.00,
    245,
    NOW() - INTERVAL '23 hours 3 minutes'
)

ON CONFLICT (attempt_id, question_index) DO UPDATE SET
    student_answer     = EXCLUDED.student_answer,
    is_correct         = EXCLUDED.is_correct,
    points_earned      = EXCLUDED.points_earned,
    max_points         = EXCLUDED.max_points,
    time_spent_seconds = EXCLUDED.time_spent_seconds,
    updated_at         = now();

COMMIT;
