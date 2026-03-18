-- =============================================================================
-- EduGo Development Seeds v2 — 004_memberships.sql
-- =============================================================================
-- 27 memberships vinculando usuarios a escuelas y unidades academicas.
--
-- Escuelas:
--   b1 = Colegio San Ignacio
--   b2 = Taller CreArte
--   b3 = Academia Global English
--
-- Mapa de memberships:
--   m001 → Carlos (U-08)    → San Ignacio → 5to A (au003)        → student
--   m002 → Carlos (U-08)    → CreArte     → Grupo Manana (au009) → student
--   m003 → Sofia (U-09)     → San Ignacio → 5to A (au003)        → student
--   m004 → Diego (U-10)     → San Ignacio → 5to B (au004)        → student
--   m005 → Valentina (U-11) → San Ignacio → 6to A (au006)        → student
--   m006 → Valentina (U-11) → Academia    → Class Monday (au014) → student
--   m007 → Mateo (U-12)     → CreArte     → Grupo Manana (au009) → student
--   m008 → Maria (U-05)     → San Ignacio → 5to A (au003)        → teacher
--   m009 → Maria (U-05)     → Academia    → Class Monday (au014) → teacher
--   m010 → Pedro (U-06)     → San Ignacio → 5to B (au004)        → teacher
--   m011 → Pedro (U-06)     → San Ignacio → 6to A (au006)        → teacher
--   m012 → Ana (U-07)       → CreArte     → Grupo Manana (au009) → teacher
--   m013 → Ana (U-07)       → CreArte     → Grupo Tarde (au011)  → teacher
--   m014 → Carmen (U-02)    → San Ignacio → NULL                 → admin
--   m015 → Roberto S (U-03) → CreArte     → NULL                 → admin
--   m016 → Lucia (U-04)     → San Ignacio → NULL                 → coordinator
--   m017 → Lucia (U-04)     → CreArte     → NULL                 → coordinator
--   m018 → Ricardo (U-13)   → San Ignacio → 5to A (au003)        → guardian
--   m019 → Ricardo (U-13)   → CreArte     → Grupo Manana (au009) → guardian
--   m020 → Patricia (U-14)  → San Ignacio → 5to A (au003)        → guardian
--   m021 → Patricia (U-14)  → San Ignacio → 5to B (au004)        → guardian
--   m022 → Miguel (U-16)    → San Ignacio → NULL                 → director
--   m023 → Laura (U-17)     → San Ignacio → NULL                 → assistant
--   m024 → Andres (U-18)    → San Ignacio → 5to A (au003)        → assistant_teacher
--   m025 → Diana (U-19)     → San Ignacio → 5to A (au003)        → observer
--   m026 → Diana (U-19)     → CreArte     → Grupo Manana (au009) → observer
--   m027 → Fernando (U-20)  → San Ignacio → 5to A (au003)        → guardian
-- =============================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- Memberships con academic_unit_id definido
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

-- Estudiantes San Ignacio
(
    'bb000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000008',   -- Carlos Mendoza
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000003',   -- 5to A
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000009',   -- Sofia Herrera
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000003',   -- 5to A
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000004',
    '00000000-0000-0000-0000-000000000010',   -- Diego Vargas
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000004',   -- 5to B
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000005',
    '00000000-0000-0000-0000-000000000011',   -- Valentina Rojas
    'b1000000-0000-0000-0000-000000000001',
    'ac000000-0000-0000-0000-000000000006',   -- 6to A
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),

-- Estudiantes multi-escuela
(
    'bb000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000008',   -- Carlos Mendoza
    'b2000000-0000-0000-0000-000000000002',   -- CreArte
    'ac000000-0000-0000-0000-000000000009',   -- Grupo Manana
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000006',
    '00000000-0000-0000-0000-000000000011',   -- Valentina Rojas
    'b3000000-0000-0000-0000-000000000003',   -- Academia
    'ac000000-0000-0000-0000-000000000014',   -- Class Monday
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000007',
    '00000000-0000-0000-0000-000000000012',   -- Mateo Fuentes
    'b2000000-0000-0000-0000-000000000002',   -- CreArte
    'ac000000-0000-0000-0000-000000000009',   -- Grupo Manana
    'student',
    '{"enrollment_type": "regular"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),

-- Docentes
(
    'bb000000-0000-0000-0000-000000000008',
    '00000000-0000-0000-0000-000000000005',   -- Maria Martinez
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000003',   -- 5to A
    'teacher',
    '{"subjects": ["Matematicas"]}'::jsonb,
    true,
    '2026-02-10 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000009',
    '00000000-0000-0000-0000-000000000005',   -- Maria Martinez
    'b3000000-0000-0000-0000-000000000003',   -- Academia
    'ac000000-0000-0000-0000-000000000014',   -- Class Monday
    'teacher',
    '{"subjects": ["English Basics A2"]}'::jsonb,
    true,
    '2026-02-10 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000010',
    '00000000-0000-0000-0000-000000000006',   -- Pedro Gonzalez
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000004',   -- 5to B
    'teacher',
    '{"subjects": ["Matematicas", "Ciencias Naturales"]}'::jsonb,
    true,
    '2026-02-10 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000011',
    '00000000-0000-0000-0000-000000000006',   -- Pedro Gonzalez
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000006',   -- 6to A
    'teacher',
    '{"subjects": ["Historia"]}'::jsonb,
    true,
    '2026-02-10 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000012',
    '00000000-0000-0000-0000-000000000007',   -- Ana Ruiz
    'b2000000-0000-0000-0000-000000000002',   -- CreArte
    'ac000000-0000-0000-0000-000000000009',   -- Grupo Manana
    'teacher',
    '{"subjects": ["Tecnicas de Pintura"]}'::jsonb,
    true,
    '2026-02-10 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000013',
    '00000000-0000-0000-0000-000000000007',   -- Ana Ruiz
    'b2000000-0000-0000-0000-000000000002',   -- CreArte
    'ac000000-0000-0000-0000-000000000011',   -- Grupo Tarde
    'teacher',
    '{"subjects": ["Fundamentos de Escultura"]}'::jsonb,
    true,
    '2026-02-10 09:00:00+00'
),

-- Assistant Teacher / Observer / Guardian (nuevos)
(
    'bb000000-0000-0000-0000-000000000024',
    '00000000-0000-0000-0000-000000000018',   -- Andres Gomez
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000003',   -- 5to A
    'assistant_teacher',
    '{"subjects": ["Matematicas"]}'::jsonb,
    true,
    '2026-02-15 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000025',
    '00000000-0000-0000-0000-000000000019',   -- Diana Lopez
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000003',   -- 5to A
    'observer',
    '{"scope": "unit"}'::jsonb,
    true,
    '2026-02-20 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000026',
    '00000000-0000-0000-0000-000000000019',   -- Diana Lopez
    'b2000000-0000-0000-0000-000000000002',   -- CreArte
    'ac000000-0000-0000-0000-000000000009',   -- Grupo Manana
    'observer',
    '{"scope": "unit"}'::jsonb,
    true,
    '2026-02-20 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000027',
    '00000000-0000-0000-0000-000000000020',   -- Fernando Ruiz
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000003',   -- 5to A
    'guardian',
    '{"ward_student_id": "00000000-0000-0000-0000-000000000008", "relationship": "tio"}'::jsonb,
    true,
    '2026-03-10 08:00:00+00'
),

-- Guardians
(
    'bb000000-0000-0000-0000-000000000018',
    '00000000-0000-0000-0000-000000000013',   -- Ricardo Mendoza
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000003',   -- 5to A
    'guardian',
    '{"ward_student_id": "00000000-0000-0000-0000-000000000008", "relationship": "padre"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000019',
    '00000000-0000-0000-0000-000000000013',   -- Ricardo Mendoza
    'b2000000-0000-0000-0000-000000000002',   -- CreArte
    'ac000000-0000-0000-0000-000000000009',   -- Grupo Manana
    'guardian',
    '{"ward_student_id": "00000000-0000-0000-0000-000000000008", "relationship": "padre"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000020',
    '00000000-0000-0000-0000-000000000014',   -- Patricia Herrera
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000003',   -- 5to A
    'guardian',
    '{"ward_student_id": "00000000-0000-0000-0000-000000000009", "relationship": "madre"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000021',
    '00000000-0000-0000-0000-000000000014',   -- Patricia Herrera
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    'ac000000-0000-0000-0000-000000000004',   -- 5to B
    'guardian',
    '{"ward_student_id": "00000000-0000-0000-0000-000000000010", "relationship": "tutora"}'::jsonb,
    true,
    '2026-03-01 08:00:00+00'
)

ON CONFLICT (user_id, school_id, academic_unit_id, role) DO UPDATE SET
    metadata   = EXCLUDED.metadata,
    is_active  = EXCLUDED.is_active,
    updated_at = now();

-- -------------------------------------------------------------------------
-- Memberships de admin/coordinator a nivel escuela (academic_unit_id = NULL)
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
    'bb000000-0000-0000-0000-000000000014',
    '00000000-0000-0000-0000-000000000002',   -- Carmen Valdes
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    NULL,
    'admin',
    '{"scope": "school"}'::jsonb,
    true,
    '2026-01-15 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000015',
    '00000000-0000-0000-0000-000000000003',   -- Roberto Silva
    'b2000000-0000-0000-0000-000000000002',   -- CreArte
    NULL,
    'admin',
    '{"scope": "school"}'::jsonb,
    true,
    '2026-01-15 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000016',
    '00000000-0000-0000-0000-000000000004',   -- Lucia Fernandez
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    NULL,
    'coordinator',
    '{"scope": "school"}'::jsonb,
    true,
    '2026-02-10 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000017',
    '00000000-0000-0000-0000-000000000004',   -- Lucia Fernandez
    'b2000000-0000-0000-0000-000000000002',   -- CreArte
    NULL,
    'coordinator',
    '{"scope": "school"}'::jsonb,
    true,
    '2026-02-10 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000022',
    '00000000-0000-0000-0000-000000000016',   -- Miguel Castillo
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    NULL,
    'director',
    '{"scope": "school"}'::jsonb,
    true,
    '2026-01-20 09:00:00+00'
),
(
    'bb000000-0000-0000-0000-000000000023',
    '00000000-0000-0000-0000-000000000017',   -- Laura Pena
    'b1000000-0000-0000-0000-000000000001',   -- San Ignacio
    NULL,
    'assistant',
    '{"scope": "school"}'::jsonb,
    true,
    '2026-02-01 09:00:00+00'
)
ON CONFLICT (id) DO UPDATE SET
    metadata   = EXCLUDED.metadata,
    is_active  = EXCLUDED.is_active,
    updated_at = now();

COMMIT;
