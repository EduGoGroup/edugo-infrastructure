-- Seeds de memberships para api-admin
-- Ejecutar después de seeds de users y academic_units
-- Demuestra diferentes roles: teacher, student, coordinator, admin

-- ==============================================================
-- Memberships para Liceo Técnico Santiago
-- ==============================================================

-- Admin de prueba (11111111-1111-1111-1111-111111111111) como ADMIN de la escuela
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active) VALUES
('a1000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', '44444444-4444-4444-4444-444444444444', 'a1000000-0000-0000-0000-000000000001', 'admin', true)
ON CONFLICT DO NOTHING;

-- Docente de prueba (22222222-2222-2222-2222-222222222222) como TEACHER en 1° Medio A
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, metadata) VALUES
('a2000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', '44444444-4444-4444-4444-444444444444', 'a1110000-0000-0000-0000-000000000001', 'teacher', true, '{"subject": "Matemáticas", "hours_per_week": 4}')
ON CONFLICT DO NOTHING;

-- Mismo docente como COORDINATOR del Club de Robótica
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, metadata) VALUES
('a2100000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', '44444444-4444-4444-4444-444444444444', 'a1400000-0000-0000-0000-000000000001', 'coordinator', true, '{"position": "Coordinador de Club"}')
ON CONFLICT DO NOTHING;

-- Estudiante de prueba (33333333-3333-3333-3333-333333333333) como STUDENT en 1° Medio A
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active) VALUES
('a3000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444', 'a1110000-0000-0000-0000-000000000001', 'student', true)
ON CONFLICT DO NOTHING;

-- Mismo estudiante en Club de Robótica
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active) VALUES
('a3100000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444', 'a1400000-0000-0000-0000-000000000001', 'student', true)
ON CONFLICT DO NOTHING;

-- Admin como TEACHER en otra clase (1° Medio B) para demostrar múltiples roles
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, metadata) VALUES
('a4000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', '44444444-4444-4444-4444-444444444444', 'a1120000-0000-0000-0000-000000000001', 'teacher', true, '{"subject": "Historia", "hours_per_week": 3}')
ON CONFLICT DO NOTHING;

-- Estudiante como ASSISTANT en Club de Ciencias
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, metadata) VALUES
('a5000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444', 'a1500000-0000-0000-0000-000000000001', 'assistant', true, '{"area": "Club de Ciencias", "schedule": "Lunes y Miércoles 15:00-17:00"}')
ON CONFLICT DO NOTHING;
