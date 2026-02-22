-- =============================================================================
-- EduGo Development Seeds — 003_users.sql
-- =============================================================================
-- Crea 13 usuarios de prueba con contraseña unificada: "EduGoTest123!"
--
-- Hash bcrypt (cost=10) de "EduGoTest123!":
--   $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a
--
-- Usuarios por categoría:
--   Sistema  : u001 → super@edugo.test (Super Admin de plataforma)
--   Admins   : u002 → admin.primaria@edugo.test
--              u003 → admin.secundario@edugo.test
--   Coord.   : u004 → coord.primaria@edugo.test
--   Docentes : u005 → teacher.math@edugo.test     (María García)
--              u006 → teacher.science@edugo.test   (Juan Martínez)
--              u007 → teacher.history@edugo.test   (Ana López)
--   Alumnos  : u008 → student.carlos@edugo.test   (Carlos González)
--              u009 → student.sofia@edugo.test     (Sofía Rodríguez)
--              u010 → student.miguel@edugo.test    (Miguel Torres)
--              u011 → student.laura@edugo.test     (Laura Martínez)
--   Tutores  : u012 → guardian.roberto@edugo.test (Roberto González)
--              u013 → guardian.patricia@edugo.test (Patricia Torres)
-- =============================================================================

BEGIN;

INSERT INTO public.users (
    id,
    email,
    password_hash,
    first_name,
    last_name,
    is_active
) VALUES

-- -------------------------------------------------------------------------
-- Sistema / Plataforma
-- -------------------------------------------------------------------------
(
    '00000000-0000-0000-0000-000000000001',
    'super@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Super',
    'Admin',
    true
),

-- -------------------------------------------------------------------------
-- Administradores de escuela
-- -------------------------------------------------------------------------
(
    '00000000-0000-0000-0000-000000000002',
    'admin.primaria@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Admin',
    'Primaria',
    true
),
(
    '00000000-0000-0000-0000-000000000003',
    'admin.secundario@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Admin',
    'Secundario',
    true
),

-- -------------------------------------------------------------------------
-- Coordinadores
-- -------------------------------------------------------------------------
(
    '00000000-0000-0000-0000-000000000004',
    'coord.primaria@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Coordinador',
    'Primaria',
    true
),

-- -------------------------------------------------------------------------
-- Docentes
-- -------------------------------------------------------------------------
(
    '00000000-0000-0000-0000-000000000005',
    'teacher.math@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'María',
    'García',
    true
),
(
    '00000000-0000-0000-0000-000000000006',
    'teacher.science@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Juan',
    'Martínez',
    true
),
(
    '00000000-0000-0000-0000-000000000007',
    'teacher.history@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Ana',
    'López',
    true
),

-- -------------------------------------------------------------------------
-- Estudiantes
-- -------------------------------------------------------------------------
(
    '00000000-0000-0000-0000-000000000008',
    'student.carlos@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Carlos',
    'González',
    true
),
(
    '00000000-0000-0000-0000-000000000009',
    'student.sofia@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Sofía',
    'Rodríguez',
    true
),
(
    '00000000-0000-0000-0000-000000000010',
    'student.miguel@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Miguel',
    'Torres',
    true
),
(
    '00000000-0000-0000-0000-000000000011',
    'student.laura@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Laura',
    'Martínez',
    true
),

-- -------------------------------------------------------------------------
-- Tutores / Apoderados
-- -------------------------------------------------------------------------
(
    '00000000-0000-0000-0000-000000000012',
    'guardian.roberto@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Roberto',
    'González',
    true
),
(
    '00000000-0000-0000-0000-000000000013',
    'guardian.patricia@edugo.test',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LPVyKmqMY.a',
    'Patricia',
    'Torres',
    true
)

ON CONFLICT (id) DO UPDATE SET
    email         = EXCLUDED.email,
    password_hash = EXCLUDED.password_hash,
    first_name    = EXCLUDED.first_name,
    last_name     = EXCLUDED.last_name,
    is_active     = EXCLUDED.is_active,
    updated_at    = now();

COMMIT;
