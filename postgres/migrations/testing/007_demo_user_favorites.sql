-- Mock Data: Favoritos de usuarios para testing
-- Materiales marcados como favoritos por estudiantes

INSERT INTO user_favorites (user_id, material_id) VALUES
-- Student1 tiene 3 favoritos
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    (SELECT id FROM materials ORDER BY created_at LIMIT 1)
),
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 1)
),
(
    (SELECT id FROM users WHERE email = 'student1@edugo.test'),
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 2)
),

-- Student2 tiene 2 favoritos
(
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    (SELECT id FROM materials ORDER BY created_at LIMIT 1)
),
(
    (SELECT id FROM users WHERE email = 'student2@edugo.test'),
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 3)
),

-- Student3 tiene 1 favorito
(
    (SELECT id FROM users WHERE email = 'student3@edugo.test'),
    (SELECT id FROM materials ORDER BY created_at LIMIT 1 OFFSET 1)
),

-- Teacher Math marca un material como favorito (referencia para sus clases)
(
    (SELECT id FROM users WHERE email = 'teacher.math@edugo.test'),
    (SELECT id FROM materials ORDER BY created_at LIMIT 1)
);
