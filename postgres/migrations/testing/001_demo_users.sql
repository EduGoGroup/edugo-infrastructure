-- ====================================================================
-- DATOS DE PRUEBA: Usuarios y asignaciones de roles para testing
-- ====================================================================
-- Contraseña para TODOS los usuarios: edugo2024
-- Hash bcrypt: $2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6
-- VERSIÓN: postgres/v0.16.5
-- ====================================================================
-- NOTA: Los roles se insertan desde seeds/002_seed_roles.sql
-- Este archivo SOLO crea usuarios y asigna roles existentes.
-- IDs de referencia de seeds:
--   super_admin:  10000000-0000-0000-0000-000000000001
--   teacher:      10000000-0000-0000-0000-000000000007
--   student:      10000000-0000-0000-0000-000000000009
--   guardian:     10000000-0000-0000-0000-000000000010
-- ====================================================================

-- ====================================================================
-- 1. INSERTAR USUARIOS DE PRUEBA
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
-- 2. ASIGNAR ROLES A USUARIOS (RBAC)
--    Referencia: IDs de roles definidos en seeds/002_seed_roles.sql
-- ====================================================================
INSERT INTO user_roles (id, user_id, role_id, school_id, academic_unit_id, is_active, granted_at) VALUES
-- Admin (rol super_admin a nivel sistema)
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '10000000-0000-0000-0000-000000000001', NULL, NULL, true, NOW()),

-- Teachers (rol teacher a nivel unit)
('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '10000000-0000-0000-0000-000000000007', NULL, NULL, true, NOW()),
('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', '10000000-0000-0000-0000-000000000007', NULL, NULL, true, NOW()),

-- Students (rol student a nivel unit)
('b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', '10000000-0000-0000-0000-000000000009', NULL, NULL, true, NOW()),
('b5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'a5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', '10000000-0000-0000-0000-000000000009', NULL, NULL, true, NOW()),
('b6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', '10000000-0000-0000-0000-000000000009', NULL, NULL, true, NOW()),

-- Guardians (rol guardian a nivel unit)
('b7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'a7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', '10000000-0000-0000-0000-000000000010', NULL, NULL, true, NOW()),
('b8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'a8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', '10000000-0000-0000-0000-000000000010', NULL, NULL, true, NOW());

-- ====================================================================
-- VERIFICACIÓN (Para logs)
-- ====================================================================
-- Se han creado:
-- - 8 usuarios con contraseña: edugo2024
-- - 8 asignaciones de roles (user_roles) usando roles de seeds
