-- ====================================================================
-- DATOS DE PRUEBA: Usuarios y Roles de demostración para testing
-- ====================================================================
-- Contraseña para TODOS los usuarios: edugo2024
-- Hash bcrypt: $2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6
-- VERSIÓN: postgres/v0.16.3
-- ====================================================================

-- ====================================================================
-- 1. INSERTAR ROLES BÁSICOS
-- ====================================================================
INSERT INTO roles (id, name, display_name, description, scope, is_active) VALUES
('11111111-1111-4111-a111-111111111111', 'admin', 'Administrador', 'Administrador del sistema con acceso completo', 'system', true),
('22222222-2222-4222-a222-222222222222', 'teacher', 'Docente', 'Profesor con acceso a gestión de cursos y evaluaciones', 'school', true),
('33333333-3333-4333-a333-333333333333', 'student', 'Estudiante', 'Estudiante con acceso a materiales y evaluaciones', 'school', true),
('44444444-4444-4444-a444-444444444444', 'guardian', 'Apoderado', 'Apoderado con acceso a información de sus hijos', 'school', true);

-- ====================================================================
-- 2. INSERTAR USUARIOS DE PRUEBA
-- ====================================================================
INSERT INTO users (id, email, password_hash, first_name, last_name, is_active, created_at, updated_at) VALUES
-- Admin (email: admin@edugo.test, password: edugo2024)
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'Admin', 'Demo', true, NOW(), NOW()),

-- Teachers (password: edugo2024)
('a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'teacher.math@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'María', 'García', true, NOW(), NOW()),
('a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'teacher.science@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'Juan', 'Pérez', true, NOW(), NOW()),

-- Students (password: edugo2024)
('a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'student1@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'Carlos', 'Rodríguez', true, NOW(), NOW()),
('a5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'student2@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'Ana', 'Martínez', true, NOW(), NOW()),
('a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'student3@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'Luis', 'González', true, NOW(), NOW()),

-- Guardians (password: edugo2024)
('a7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'guardian1@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'Roberto', 'Fernández', true, NOW(), NOW()),
('a8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'guardian2@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'Patricia', 'López', true, NOW(), NOW());

-- ====================================================================
-- 3. ASIGNAR ROLES A USUARIOS (RBAC)
-- ====================================================================
INSERT INTO user_roles (id, user_id, role_id, school_id, academic_unit_id, is_active, granted_at) VALUES
-- Admin (rol a nivel sistema)
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '11111111-1111-4111-a111-111111111111', NULL, NULL, true, NOW()),

-- Teachers (roles a nivel escuela - sin school_id específico por ahora)
('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '22222222-2222-4222-a222-222222222222', NULL, NULL, true, NOW()),
('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '22222222-2222-4222-a222-222222222222', NULL, NULL, true, NOW()),

-- Students (roles a nivel escuela - sin school_id específico por ahora)
('b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', '33333333-3333-4333-a333-333333333333', NULL, NULL, true, NOW()),
('b5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'a5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', '33333333-3333-4333-a333-333333333333', NULL, NULL, true, NOW()),
('b6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', '33333333-3333-4333-a333-333333333333', NULL, NULL, true, NOW()),

-- Guardians (roles a nivel escuela - sin school_id específico por ahora)
('b7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'a7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', '44444444-4444-4444-a444-444444444444', NULL, NULL, true, NOW()),
('b8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'a8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', '44444444-4444-4444-a444-444444444444', NULL, NULL, true, NOW());

-- ====================================================================
-- VERIFICACIÓN (Para logs)
-- ====================================================================
-- Se han creado:
-- - 4 roles: admin, teacher, student, guardian
-- - 8 usuarios con contraseña: edugo2024
-- - 8 asignaciones de roles (user_roles)
