-- =============================================================================
-- EduGo Development Seeds — 004_memberships.sql
-- =============================================================================
-- Crea 12 memberships que vinculan usuarios a escuelas y unidades académicas.
--
-- La constraint unique es: (user_id, school_id, academic_unit_id, role)
-- Para admins y coordinadores con academic_unit_id=NULL se debe tener cuidado:
-- PostgreSQL trata NULL como distinto en constraints unique, por eso los
-- memberships de admin sin unidad usan ON CONFLICT DO NOTHING (el conflict
-- se dispararía solo si existe un registro idéntico incluyendo NULL=NULL,
-- que en Postgres no ocurre). Se usa DO NOTHING como comportamiento seguro.
--
-- Mapa de memberships:
--   m001 → carlos  (u008) → primaria (b1) → Clase 1-A (au003) → student
--   m002 → sofia   (u009) → primaria (b1) → Clase 1-A (au003) → student
--   m003 → miguel  (u010) → primaria (b1) → Clase 1-B (au004) → student
--   m004 → math    (u005) → primaria (b1) → Clase 1-A (au003) → teacher
--   m005 → science (u006) → primaria (b1) → Clase 1-B (au004) → teacher
--   m006 → laura   (u011) → secundario (b2) → Clase 10-A (au009) → student
--   m007 → history (u007) → secundario (b2) → Clase 10-A (au009) → teacher
--   m008 → admin.primaria (u002) → primaria (b1) → NULL → admin
--   m009 → admin.secundario (u003) → secundario (b2) → NULL → admin
--   m010 → coord.primaria (u004) → primaria (b1) → Primer Grado (au002) → coordinator
--   m011 → guardian.roberto (u012) → primaria (b1) → Clase 1-A (au003) → guardian
--   m012 → guardian.patricia (u013) → primaria (b1) → Clase 1-B (au004) → guardian
-- =============================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- Memberships con academic_unit_id definido
-- (conflict detectable por la constraint unique compuesta)
-- -------------------------------------------------------------------------
INSERT INTO academic.memberships (
    id,
    user_id,
    school_id,
    academic_unit_id,
    role,
    metadata,
    is_active,
    enrolled_at
) VALUES

-- Estudiantes Clase 1-A (Escuela Primaria)
(
    'bb000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000008',   -- carlos
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    'ac000000-0000-0000-0000-000000000003',   -- Clase 1-A
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2024-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000009',   -- sofia
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000003',   -- Clase 1-A
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2024-03-01 08:00:00+00'
),

-- Estudiante Clase 1-B (Escuela Primaria)
(
    'bb000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000010',   -- miguel
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000004',   -- Clase 1-B
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2024-03-01 08:00:00+00'
),

-- Docente Matemáticas → Clase 1-A
(
    'bb000000-0000-0000-0000-000000000004',
    '00000000-0000-0000-0000-000000000005',   -- teacher.math (María García)
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000003',   -- Clase 1-A
    'teacher',
    '{"subjects": ["Matemáticas"]}'::jsonb,
    true,
    '2024-02-15 09:00:00+00'
),

-- Docente Ciencias → Clase 1-B
(
    'bb000000-0000-0000-0000-000000000005',
    '00000000-0000-0000-0000-000000000006',   -- teacher.science (Juan Martínez)
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000004',   -- Clase 1-B
    'teacher',
    '{"subjects": ["Ciencias Naturales"]}'::jsonb,
    true,
    '2024-02-15 09:00:00+00'
),

-- Estudiante Clase 10-A (Colegio Secundario)
(
    'bb000000-0000-0000-0000-000000000006',
    '00000000-0000-0000-0000-000000000011',   -- laura
    'b2000000-0000-0000-0000-000000000002',   -- Colegio Secundario Demo
    'ac000000-0000-0000-0000-000000000009',   -- Clase 10-A
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2024-03-01 08:00:00+00'
),

-- Docente Historia → Clase 10-A
(
    'bb000000-0000-0000-0000-000000000007',
    '00000000-0000-0000-0000-000000000007',   -- teacher.history (Ana López)
    'b2000000-0000-0000-0000-000000000002',
    'ac000000-0000-0000-0000-000000000009',   -- Clase 10-A
    'teacher',
    '{"subjects": ["Historia"]}'::jsonb,
    true,
    '2024-02-15 09:00:00+00'
),

-- teacher.math ALSO as coordinator at Colegio Secundario (dual-school role test)
(
    'bb000000-0000-0000-0000-000000000008',
    '00000000-0000-0000-0000-000000000005',   -- teacher.math (María García)
    'b2000000-0000-0000-0000-000000000002',   -- Colegio Secundario Demo
    NULL,                                      -- school-level coordinator
    'coordinator',
    '{"scope": "school"}'::jsonb,
    true,
    '2024-02-15 09:00:00+00'
),

-- Coordinadora → Primer Grado (Escuela Primaria)
(
    'bb000000-0000-0000-0000-000000000010',
    '00000000-0000-0000-0000-000000000004',   -- coord.primaria
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000002',   -- Primer Grado
    'coordinator',
    '{"scope": "grade"}'::jsonb,
    true,
    '2024-02-10 09:00:00+00'
),

-- Tutores → Clase 1-A y 1-B
(
    'bb000000-0000-0000-0000-000000000011',
    '00000000-0000-0000-0000-000000000012',   -- guardian.roberto (padre de Carlos)
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000003',   -- Clase 1-A
    'guardian',
    '{"ward_student_id": "00000000-0000-0000-0000-000000000008", "relationship": "padre"}'::jsonb,
    true,
    '2024-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000012',
    '00000000-0000-0000-0000-000000000013',   -- guardian.patricia (madre de Miguel)
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000004',   -- Clase 1-B
    'guardian',
    '{"ward_student_id": "00000000-0000-0000-0000-000000000010", "relationship": "madre"}'::jsonb,
    true,
    '2024-03-01 08:00:00+00'
)

ON CONFLICT (user_id, school_id, academic_unit_id, role) DO UPDATE SET
    metadata   = EXCLUDED.metadata,
    is_active  = EXCLUDED.is_active,
    updated_at = now();

-- -------------------------------------------------------------------------
-- Memberships de admin a nivel de escuela completa (academic_unit_id = NULL)
-- NULL no satisface la constraint unique, por eso se usa INSERT separado
-- con ON CONFLICT DO NOTHING (idempotente por id).
-- -------------------------------------------------------------------------
INSERT INTO academic.memberships (
    id,
    user_id,
    school_id,
    academic_unit_id,
    role,
    metadata,
    is_active,
    enrolled_at
) VALUES
(
    'bb000000-0000-0000-0000-000000000008',
    '00000000-0000-0000-0000-000000000002',   -- admin.primaria
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    NULL,                                      -- scope: toda la escuela
    'admin',
    '{"scope": "school"}'::jsonb,
    true,
    '2024-01-15 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000009',
    '00000000-0000-0000-0000-000000000003',   -- admin.secundario
    'b2000000-0000-0000-0000-000000000002',   -- Colegio Secundario Demo
    NULL,
    'admin',
    '{"scope": "school"}'::jsonb,
    true,
    '2024-01-15 09:00:00+00'
)
ON CONFLICT (id) DO UPDATE SET
    metadata   = EXCLUDED.metadata,
    is_active  = EXCLUDED.is_active,
    updated_at = now();

COMMIT;
