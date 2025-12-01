-- Mock Data: Contextos activos de usuarios para testing
-- Establece la escuela activa para usuarios demo

INSERT INTO user_active_context (user_id, school_id, unit_id) VALUES
-- Admin tiene como activa la primera escuela, primera unidad
(
    (SELECT id FROM users WHERE email = 'admin@edugo.test'),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    (SELECT id FROM academic_units ORDER BY created_at LIMIT 1)
),

-- Teacher Math tiene su escuela asignada
(
    (SELECT id FROM users WHERE email = 'teacher.math@edugo.test'),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    (SELECT id FROM academic_units ORDER BY created_at LIMIT 1)
),

-- Teacher Science tiene su escuela asignada
(
    (SELECT id FROM users WHERE email = 'teacher.science@edugo.test'),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    (SELECT id FROM academic_units ORDER BY created_at LIMIT 1 OFFSET 1)
),

-- Student1 tiene contexto en primera escuela, primera unidad
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    (SELECT id FROM academic_units ORDER BY created_at LIMIT 1)
),

-- Student2 tiene contexto en primera escuela, primera unidad
(
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    (SELECT id FROM academic_units ORDER BY created_at LIMIT 1)
),

-- Student3 tiene contexto en primera escuela, segunda unidad
(
    (SELECT id FROM users WHERE email = 'student3@edugo.test'),
    (SELECT id FROM schools ORDER BY created_at LIMIT 1),
    (SELECT id FROM academic_units ORDER BY created_at LIMIT 1 OFFSET 1)
);
