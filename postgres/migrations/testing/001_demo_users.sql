-- Mock Data: Usuarios de demostración para testing
-- Contraseña para TODOS los usuarios: edugo2024

INSERT INTO users (id, email, password_hash, role, first_name, last_name, created_at, updated_at) VALUES
-- Admin (email: admin@edugo.test, password: edugo2024)
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'admin', 'Admin', 'Demo', NOW(), NOW()),

-- Teachers (password: edugo2024)
('a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'teacher.math@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'teacher', 'María', 'García', NOW(), NOW()),
('a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'teacher.science@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'teacher', 'Juan', 'Pérez', NOW(), NOW()),

-- Students (password: edugo2024)
('a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'student1@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'student', 'Carlos', 'Rodríguez', NOW(), NOW()),
('a5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'student2@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'student', 'Ana', 'Martínez', NOW(), NOW()),
('a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'student3@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'student', 'Luis', 'González', NOW(), NOW()),

-- Guardians (password: edugo2024)
('a7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'guardian1@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'guardian', 'Roberto', 'Fernández', NOW(), NOW()),
('a8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'guardian2@edugo.test', '$2a$10$x0lpvYBLh8dCiMYskYzD1.y2TfeXcQh7QbBXIO5Xepi3SIgC2FtY6', 'guardian', 'Patricia', 'López', NOW(), NOW());
