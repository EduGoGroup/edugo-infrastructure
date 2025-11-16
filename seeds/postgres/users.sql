-- Seeds de usuarios de prueba
-- Ejecutar después de migración 001_create_users

INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, email_verified) VALUES
-- Admin de prueba
('11111111-1111-1111-1111-111111111111', 'admin@edugo.com', '$2a$10$hash', 'Admin', 'Sistema', 'admin', true, true),

-- Docente de prueba
('22222222-2222-2222-2222-222222222222', 'teacher@edugo.com', '$2a$10$hash', 'María', 'González', 'teacher', true, true),

-- Estudiante de prueba
('33333333-3333-3333-3333-333333333333', 'student@edugo.com', '$2a$10$hash', 'Juan', 'Pérez', 'student', true, true)

ON CONFLICT (email) DO NOTHING;
