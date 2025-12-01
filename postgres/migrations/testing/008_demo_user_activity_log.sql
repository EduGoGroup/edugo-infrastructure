-- Mock Data: Log de actividades de usuarios para testing
-- Simula actividad reciente de estudiantes

INSERT INTO user_activity_log (user_id, activity_type, material_id, school_id, metadata, created_at) VALUES
-- Student1: Flujo completo de estudio
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    'material_started',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"page": 1}'::jsonb,
    NOW() - INTERVAL '2 hours'
),
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    'material_progress',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"page": 5, "total_pages": 10, "time_spent_seconds": 300}'::jsonb,
    NOW() - INTERVAL '1 hour 45 minutes'
),
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    'material_completed',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"total_time_seconds": 3600}'::jsonb,
    NOW() - INTERVAL '1 hour'
),
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    'summary_viewed',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"summary_length_chars": 500, "read_time_seconds": 45}'::jsonb,
    NOW() - INTERVAL '50 minutes'
),
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    'quiz_started',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"question_count": 10}'::jsonb,
    NOW() - INTERVAL '30 minutes'
),
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    'quiz_passed',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"score": 90, "correct_answers": 9, "total_questions": 10, "time_seconds": 600}'::jsonb,
    NOW() - INTERVAL '20 minutes'
),

-- Student2: Actividad reciente
(
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    'material_started',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"page": 1}'::jsonb,
    NOW() - INTERVAL '3 hours'
),
(
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    'material_progress',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"page": 3, "total_pages": 8, "time_spent_seconds": 180}'::jsonb,
    NOW() - INTERVAL '2 hours 30 minutes'
),
(
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    'quiz_started',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"question_count": 5}'::jsonb,
    NOW() - INTERVAL '1 hour'
),
(
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    'quiz_failed',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 1),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"score": 40, "correct_answers": 2, "total_questions": 5, "time_seconds": 300}'::jsonb,
    NOW() - INTERVAL '50 minutes'
),

-- Student3: Pocas actividades
(
    (SELECT id FROM users WHERE email = 'student3@edugo.test'),
    'material_started',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 2),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"page": 1}'::jsonb,
    NOW() - INTERVAL '5 hours'
),
(
    (SELECT id FROM users WHERE email = 'student3@edugo.test'),
    'material_progress',
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 2),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    '{"page": 2, "total_pages": 6, "time_spent_seconds": 120}'::jsonb,
    NOW() - INTERVAL '4 hours 30 minutes'
);
