-- =============================================================================
-- EduGo Development Seeds — 005_user_roles.sql
-- =============================================================================
-- Asigna roles RBAC a los 13 usuarios de desarrollo.
-- Implementa R4: school_id en user_roles para aislar permisos por escuela.
--
-- Reglas de scoping:
--   - super_admin    → school_id=NULL  (acceso global a la plataforma)
--   - school_admin   → school_id=bX   (acceso solo a su escuela)
--   - coordinator    → school_id=bX   (acceso a su escuela)
--   - teacher        → school_id=bX   (acceso a su escuela)
--   - student        → school_id=bX   (acceso a su escuela)
--   - guardian       → school_id=bX   (acceso a su escuela)
--
-- IDs de roles (fijos, seeds de producción):
--   10000000-...0001 → super_admin
--   10000000-...0003 → school_admin
--   10000000-...0005 → school_coordinator
--   10000000-...0007 → teacher
--   10000000-...0009 → student
--   10000000-...0010 → guardian
--
-- La constraint unique es: (user_id, role_id, school_id, academic_unit_id)
-- NULL en school_id/academic_unit_id se trata como valor único en Postgres
-- para constraints unique compuestas (NULL != NULL), por eso el super_admin
-- usa INSERT separado con ON CONFLICT por id.
-- =============================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- Super Admin — sin scope de escuela (acceso global)
-- -------------------------------------------------------------------------
INSERT INTO iam.user_roles (
    id,
    user_id,
    role_id,
    school_id,
    academic_unit_id,
    is_active,
    granted_by,
    granted_at
) VALUES (
    'cc000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',   -- super@edugo.test
    '10000000-0000-0000-0000-000000000001',   -- super_admin
    NULL,
    NULL,
    true,
    '00000000-0000-0000-0000-000000000001',   -- auto-asignado
    '2024-01-01 00:00:00'
)
ON CONFLICT (id) DO UPDATE SET
    is_active  = EXCLUDED.is_active,
    updated_at = CURRENT_TIMESTAMP;

-- -------------------------------------------------------------------------
-- Roles con scope de escuela (school_id definido)
-- La constraint unique funciona correctamente cuando school_id != NULL.
-- -------------------------------------------------------------------------
INSERT INTO iam.user_roles (
    id,
    user_id,
    role_id,
    school_id,
    academic_unit_id,
    is_active,
    granted_by,
    granted_at
) VALUES

-- Administradores de escuela
(
    'cc000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000002',   -- admin.primaria
    '10000000-0000-0000-0000-000000000003',   -- school_admin
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000001',   -- granted_by super
    '2024-01-15 09:00:00'
),
(
    'cc000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000003',   -- admin.secundario
    '10000000-0000-0000-0000-000000000003',   -- school_admin
    'b2000000-0000-0000-0000-000000000002',   -- Colegio Secundario Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000001',   -- granted_by super
    '2024-01-15 09:00:00'
),

-- Coordinadora de grado
(
    'cc000000-0000-0000-0000-000000000004',
    '00000000-0000-0000-0000-000000000004',   -- coord.primaria
    '10000000-0000-0000-0000-000000000005',   -- school_coordinator
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000002',   -- granted_by admin.primaria
    '2024-02-01 09:00:00'
),

-- Docentes
(
    'cc000000-0000-0000-0000-000000000005',
    '00000000-0000-0000-0000-000000000005',   -- teacher.math (María García)
    '10000000-0000-0000-0000-000000000007',   -- teacher
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000002',   -- granted_by admin.primaria
    '2024-02-10 09:00:00'
),
(
    'cc000000-0000-0000-0000-000000000006',
    '00000000-0000-0000-0000-000000000006',   -- teacher.science (Juan Martínez)
    '10000000-0000-0000-0000-000000000007',   -- teacher
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000002',
    '2024-02-10 09:00:00'
),
(
    'cc000000-0000-0000-0000-000000000007',
    '00000000-0000-0000-0000-000000000007',   -- teacher.history (Ana López)
    '10000000-0000-0000-0000-000000000007',   -- teacher
    'b2000000-0000-0000-0000-000000000002',   -- Colegio Secundario Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000003',   -- granted_by admin.secundario
    '2024-02-10 09:00:00'
),
-- teacher.math ALSO as coordinator at Colegio Secundario (dual-school role test)
(
    'cc000000-0000-0000-0000-000000000020',
    '00000000-0000-0000-0000-000000000005',   -- teacher.math (María García)
    '10000000-0000-0000-0000-000000000005',   -- school_coordinator
    'b2000000-0000-0000-0000-000000000002',   -- Colegio Secundario Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000003',   -- granted_by admin.secundario
    '2024-02-15 09:00:00'
),

-- Estudiantes
(
    'cc000000-0000-0000-0000-000000000008',
    '00000000-0000-0000-0000-000000000008',   -- carlos
    '10000000-0000-0000-0000-000000000009',   -- student
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000002',
    '2024-03-01 08:00:00'
),
(
    'cc000000-0000-0000-0000-000000000009',
    '00000000-0000-0000-0000-000000000009',   -- sofia
    '10000000-0000-0000-0000-000000000009',   -- student
    'b1000000-0000-0000-0000-000000000001',
    NULL,
    true,
    '00000000-0000-0000-0000-000000000002',
    '2024-03-01 08:00:00'
),
(
    'cc000000-0000-0000-0000-000000000010',
    '00000000-0000-0000-0000-000000000010',   -- miguel
    '10000000-0000-0000-0000-000000000009',   -- student
    'b1000000-0000-0000-0000-000000000001',
    NULL,
    true,
    '00000000-0000-0000-0000-000000000002',
    '2024-03-01 08:00:00'
),
(
    'cc000000-0000-0000-0000-000000000011',
    '00000000-0000-0000-0000-000000000011',   -- laura
    '10000000-0000-0000-0000-000000000009',   -- student
    'b2000000-0000-0000-0000-000000000002',   -- Colegio Secundario Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000003',
    '2024-03-01 08:00:00'
),

-- Tutores / Apoderados
(
    'cc000000-0000-0000-0000-000000000012',
    '00000000-0000-0000-0000-000000000012',   -- guardian.roberto
    '10000000-0000-0000-0000-000000000010',   -- guardian
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    NULL,
    true,
    '00000000-0000-0000-0000-000000000002',
    '2024-03-01 08:00:00'
),
(
    'cc000000-0000-0000-0000-000000000013',
    '00000000-0000-0000-0000-000000000013',   -- guardian.patricia
    '10000000-0000-0000-0000-000000000010',   -- guardian
    'b1000000-0000-0000-0000-000000000001',
    NULL,
    true,
    '00000000-0000-0000-0000-000000000002',
    '2024-03-01 08:00:00'
)

ON CONFLICT (user_id, role_id, school_id, academic_unit_id) DO UPDATE SET
    is_active  = EXCLUDED.is_active,
    granted_by = EXCLUDED.granted_by,
    updated_at = CURRENT_TIMESTAMP;

COMMIT;
