-- ============================================================
-- SEED: roles (11 registros)
-- Fecha: 2026-02-22
-- Fuente: Datos reales de producción (Neon)
-- Idempotente: usa ON CONFLICT DO NOTHING
-- ============================================================

INSERT INTO roles (id, name, display_name, description, scope, is_active)
VALUES
  -- Roles de sistema
  ('10000000-0000-0000-0000-000000000001', 'super_admin',        'Super Administrador',         'Administrador con acceso total al sistema',                  'system', TRUE),
  ('10000000-0000-0000-0000-000000000002', 'platform_admin',     'Administrador de Plataforma', 'Administrador de plataforma con permisos de gestión global', 'system', TRUE),

  -- Roles de escuela
  ('10000000-0000-0000-0000-000000000003', 'school_admin',       'Administrador de Escuela',    'Administrador con control total de la institución',          'school', TRUE),
  ('10000000-0000-0000-0000-000000000004', 'school_director',    'Director',                    'Director de la institución educativa',                       'school', TRUE),
  ('10000000-0000-0000-0000-000000000005', 'school_coordinator', 'Coordinador',                 'Coordinador académico de la institución',                    'school', TRUE),
  ('10000000-0000-0000-0000-000000000006', 'school_assistant',   'Asistente Administrativo',    'Personal de soporte administrativo',                         'school', TRUE),

  -- Roles de unidad
  ('10000000-0000-0000-0000-000000000007', 'teacher',            'Profesor',                    'Docente con permisos de gestión de clase',                   'unit',   TRUE),
  ('10000000-0000-0000-0000-000000000008', 'assistant_teacher',  'Profesor Asistente',          'Asistente de docente',                                       'unit',   TRUE),
  ('10000000-0000-0000-0000-000000000009', 'student',            'Estudiante',                  'Alumno inscrito en la unidad',                               'unit',   TRUE),
  ('10000000-0000-0000-0000-000000000010', 'guardian',           'Apoderado',                   'Tutor legal o apoderado de estudiante',                      'unit',   TRUE),
  ('10000000-0000-0000-0000-000000000011', 'observer',           'Observador',                  'Rol de solo lectura para auditoría',                         'unit',   TRUE)
ON CONFLICT (id) DO NOTHING;
