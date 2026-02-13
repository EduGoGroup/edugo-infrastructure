-- ====================================================================
-- SEEDS: Roles predefinidos del sistema
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

INSERT INTO roles (id, name, display_name, description, scope, is_active) VALUES

-- Roles a nivel sistema
('10000000-0000-0000-0000-000000000001', 'super_admin', 'Super Administrador',
 'Administrador con acceso total al sistema', 'system', true),

('10000000-0000-0000-0000-000000000002', 'platform_admin', 'Administrador de Plataforma',
 'Administrador de plataforma con permisos de gestión global', 'system', true),

-- Roles a nivel escuela
('10000000-0000-0000-0000-000000000003', 'school_admin', 'Administrador de Escuela',
 'Administrador con control total de la institución', 'school', true),

('10000000-0000-0000-0000-000000000004', 'school_director', 'Director',
 'Director de la institución educativa', 'school', true),

('10000000-0000-0000-0000-000000000005', 'school_coordinator', 'Coordinador',
 'Coordinador académico de la institución', 'school', true),

('10000000-0000-0000-0000-000000000006', 'school_assistant', 'Asistente Administrativo',
 'Personal de soporte administrativo', 'school', true),

-- Roles a nivel unidad académica
('10000000-0000-0000-0000-000000000007', 'teacher', 'Profesor',
 'Docente con permisos de gestión de clase', 'unit', true),

('10000000-0000-0000-0000-000000000008', 'assistant_teacher', 'Profesor Asistente',
 'Asistente de docente', 'unit', true),

('10000000-0000-0000-0000-000000000009', 'student', 'Estudiante',
 'Alumno inscrito en la unidad', 'unit', true),

('10000000-0000-0000-0000-000000000010', 'guardian', 'Apoderado',
 'Tutor legal o apoderado de estudiante', 'unit', true),

('10000000-0000-0000-0000-000000000011', 'observer', 'Observador',
 'Rol de solo lectura para auditoría', 'unit', true);
