-- ====================================================================
-- DATOS DE PRUEBA: Usuarios y Roles de demostración para testing
-- ====================================================================
-- Contraseña para TODOS los usuarios: edugo2024
-- Hash bcrypt: $2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6
-- VERSIÓN: postgres/v0.16.2
-- ====================================================================

-- ====================================================================
-- 1. INSERTAR ROLES BÁSICOS
-- ====================================================================
INSERT INTO roles (id, name, display_name, description, scope, is_active) VALUES
('r1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin', 'Administrador', 'Administrador del sistema con acceso completo', 'system', true),
('r2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'teacher', 'Docente', 'Profesor con acceso a gestión de cursos y evaluaciones', 'school', true),
('r3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'student', 'Estudiante', 'Estudiante con acceso a materiales y evaluaciones', 'school', true),
('r4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'guardian', 'Apoderado', 'Apoderado con acceso a información de sus hijos', 'school', true);

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
('ur1ebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'r1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NULL, NULL, true, NOW()),

-- Teachers (roles a nivel escuela - sin school_id específico por ahora)
('ur2ebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'r2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', NULL, NULL, true, NOW()),
('ur3ebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'r2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', NULL, NULL, true, NOW()),

-- Students (roles a nivel escuela - sin school_id específico por ahora)
('ur4ebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'r3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', NULL, NULL, true, NOW()),
('ur5ebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'a5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'r3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', NULL, NULL, true, NOW()),
('ur6ebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'r3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', NULL, NULL, true, NOW()),

-- Guardians (roles a nivel escuela - sin school_id específico por ahora)
('ur7ebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'a7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'r4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', NULL, NULL, true, NOW()),
('ur8ebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'a8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'r4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', NULL, NULL, true, NOW());

-- ====================================================================
-- VERIFICACIÓN (Para logs)
-- ====================================================================
-- Se han creado:
-- - 4 roles: admin, teacher, student, guardian
-- - 8 usuarios con contraseña: edugo2024
-- - 8 asignaciones de roles (user_roles)
