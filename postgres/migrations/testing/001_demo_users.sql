-- Mock Data: Usuarios de demostración para testing

INSERT INTO users (id, email, password_hash, role, first_name, last_name, created_at, updated_at) VALUES
-- Admin
('usr_admin_demo', 'admin@edugo.test', '$2a$10$YourHashHere', 'admin', 'Admin', 'Demo', NOW(), NOW()),

-- Teachers
('usr_teacher_math', 'teacher.math@edugo.test', '$2a$10$YourHashHere', 'teacher', 'María', 'García', NOW(), NOW()),
('usr_teacher_science', 'teacher.science@edugo.test', '$2a$10$YourHashHere', 'teacher', 'Juan', 'Pérez', NOW(), NOW()),

-- Students
('usr_student_1', 'student1@edugo.test', '$2a$10$YourHashHere', 'student', 'Carlos', 'Rodríguez', NOW(), NOW()),
('usr_student_2', 'student2@edugo.test', '$2a$10$YourHashHere', 'student', 'Ana', 'Martínez', NOW(), NOW()),
('usr_student_3', 'student3@edugo.test', '$2a$10$YourHashHere', 'student', 'Luis', 'González', NOW(), NOW()),

-- Guardians
('usr_guardian_1', 'guardian1@edugo.test', '$2a$10$YourHashHere', 'guardian', 'Roberto', 'Fernández', NOW(), NOW()),
('usr_guardian_2', 'guardian2@edugo.test', '$2a$10$YourHashHere', 'guardian', 'Patricia', 'López', NOW(), NOW());
