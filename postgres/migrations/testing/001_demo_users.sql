-- Mock Data: Usuarios de demostración para testing

INSERT INTO users (id, email, password_hash, role, first_name, last_name, created_at, updated_at) VALUES
-- Admin
('a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@edugo.test', '$2a$10$YourHashHere', 'admin', 'Admin', 'Demo', NOW(), NOW()),

-- Teachers
('a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'teacher.math@edugo.test', '$2a$10$YourHashHere', 'teacher', 'María', 'García', NOW(), NOW()),
('a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'teacher.science@edugo.test', '$2a$10$YourHashHere', 'teacher', 'Juan', 'Pérez', NOW(), NOW()),

-- Students
('a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'student1@edugo.test', '$2a$10$YourHashHere', 'student', 'Carlos', 'Rodríguez', NOW(), NOW()),
('a5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'student2@edugo.test', '$2a$10$YourHashHere', 'student', 'Ana', 'Martínez', NOW(), NOW()),
('a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'student3@edugo.test', '$2a$10$YourHashHere', 'student', 'Luis', 'González', NOW(), NOW()),

-- Guardians
('a7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'guardian1@edugo.test', '$2a$10$YourHashHere', 'guardian', 'Roberto', 'Fernández', NOW(), NOW()),
('a8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'guardian2@edugo.test', '$2a$10$YourHashHere', 'guardian', 'Patricia', 'López', NOW(), NOW());
