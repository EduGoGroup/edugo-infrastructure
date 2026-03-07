-- =============================================================================
-- EduGo Development Seeds v2 — 002_academic_units.sql
-- =============================================================================
-- 16 unidades academicas en jerarquia (padres antes que hijos).
--
-- Colegio San Ignacio (b1000...001):
--   au001 school  → "Colegio San Ignacio"     (raiz)
--   au002 grade   → "5to Basico"              (parent=au001)
--   au003 class   → "5to A"                   (parent=au002)
--   au004 class   → "5to B"                   (parent=au002)
--   au005 grade   → "6to Basico"              (parent=au001)
--   au006 class   → "6to A"                   (parent=au005)
--
-- Taller CreArte (b2000...002):
--   au007 school  → "Taller CreArte"          (raiz)
--   au008 grade   → "Modulo Pintura"          (parent=au007)
--   au009 class   → "Grupo Manana"            (parent=au008)
--   au010 grade   → "Modulo Escultura"        (parent=au007)
--   au011 class   → "Grupo Tarde"             (parent=au010)
--
-- Academia Global English (b3000...003):
--   au012 school  → "Academia Global English" (raiz)
--   au013 grade   → "Level A2"                (parent=au012)
--   au014 class   → "Class Monday"            (parent=au013)
--   au015 grade   → "Level B1"                (parent=au012)
--   au016 class   → "Class Tuesday"           (parent=au015)
-- =============================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- Nivel 0: Unidades raiz (type=school, parent=NULL)
-- -------------------------------------------------------------------------
INSERT INTO academic.academic_units (
    id, parent_unit_id, school_id, name, code, type, description, level, academic_year, metadata, is_active
) VALUES
(
    'ac000000-0000-0000-0000-000000000001',
    NULL,
    'b1000000-0000-0000-0000-000000000001',
    'Colegio San Ignacio',
    'CSI-ROOT',
    'school',
    'Unidad raiz del Colegio San Ignacio',
    'secondary',
    0,
    '{"is_root": true}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000007',
    NULL,
    'b2000000-0000-0000-0000-000000000002',
    'Taller CreArte',
    'TCA-ROOT',
    'school',
    'Unidad raiz del Taller CreArte',
    'workshop',
    0,
    '{"is_root": true}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000012',
    NULL,
    'b3000000-0000-0000-0000-000000000003',
    'Academia Global English',
    'AGE-ROOT',
    'school',
    'Unidad raiz de la Academia Global English',
    'language',
    0,
    '{"is_root": true}'::jsonb,
    true
)
ON CONFLICT (school_id, code, academic_year) DO UPDATE SET
    name           = EXCLUDED.name,
    parent_unit_id = EXCLUDED.parent_unit_id,
    description    = EXCLUDED.description,
    level          = EXCLUDED.level,
    metadata       = EXCLUDED.metadata,
    is_active      = EXCLUDED.is_active,
    updated_at     = now();

-- -------------------------------------------------------------------------
-- Nivel 1: Grados / Modulos / Levels (type=grade, parent=school)
-- -------------------------------------------------------------------------
INSERT INTO academic.academic_units (
    id, parent_unit_id, school_id, name, code, type, description, level, academic_year, metadata, is_active
) VALUES
-- San Ignacio
(
    'ac000000-0000-0000-0000-000000000002',
    'ac000000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',
    '5to Basico',
    'GRADE-05',
    'grade',
    'Quinto ano de educacion basica, 2026',
    'secondary',
    2026,
    '{"grade_number": 5}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000005',
    'ac000000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',
    '6to Basico',
    'GRADE-06',
    'grade',
    'Sexto ano de educacion basica, 2026',
    'secondary',
    2026,
    '{"grade_number": 6}'::jsonb,
    true
),
-- CreArte
(
    'ac000000-0000-0000-0000-000000000008',
    'ac000000-0000-0000-0000-000000000007',
    'b2000000-0000-0000-0000-000000000002',
    'Modulo Pintura',
    'MOD-PINT',
    'grade',
    'Modulo de tecnicas de pintura',
    'workshop',
    2026,
    '{"module_type": "pintura"}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000010',
    'ac000000-0000-0000-0000-000000000007',
    'b2000000-0000-0000-0000-000000000002',
    'Modulo Escultura',
    'MOD-ESCL',
    'grade',
    'Modulo de fundamentos de escultura',
    'workshop',
    2026,
    '{"module_type": "escultura"}'::jsonb,
    true
),
-- Academia
(
    'ac000000-0000-0000-0000-000000000013',
    'ac000000-0000-0000-0000-000000000012',
    'b3000000-0000-0000-0000-000000000003',
    'Level A2',
    'LVL-A2',
    'grade',
    'Elementary level A2',
    'language',
    2026,
    '{"cefr_level": "A2"}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000015',
    'ac000000-0000-0000-0000-000000000012',
    'b3000000-0000-0000-0000-000000000003',
    'Level B1',
    'LVL-B1',
    'grade',
    'Intermediate level B1',
    'language',
    2026,
    '{"cefr_level": "B1"}'::jsonb,
    true
)
ON CONFLICT (school_id, code, academic_year) DO UPDATE SET
    name           = EXCLUDED.name,
    parent_unit_id = EXCLUDED.parent_unit_id,
    description    = EXCLUDED.description,
    level          = EXCLUDED.level,
    metadata       = EXCLUDED.metadata,
    is_active      = EXCLUDED.is_active,
    updated_at     = now();

-- -------------------------------------------------------------------------
-- Nivel 2: Clases / Grupos / Classes (type=class, parent=grade)
-- -------------------------------------------------------------------------
INSERT INTO academic.academic_units (
    id, parent_unit_id, school_id, name, code, type, description, level, academic_year, metadata, is_active
) VALUES
-- San Ignacio
(
    'ac000000-0000-0000-0000-000000000003',
    'ac000000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000001',
    '5to A',
    '5A',
    'class',
    'Seccion A del 5to Basico, 2026',
    'secondary',
    2026,
    '{"section": "A", "grade_number": 5}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000004',
    'ac000000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000001',
    '5to B',
    '5B',
    'class',
    'Seccion B del 5to Basico, 2026',
    'secondary',
    2026,
    '{"section": "B", "grade_number": 5}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000006',
    'ac000000-0000-0000-0000-000000000005',
    'b1000000-0000-0000-0000-000000000001',
    '6to A',
    '6A',
    'class',
    'Seccion A del 6to Basico, 2026',
    'secondary',
    2026,
    '{"section": "A", "grade_number": 6}'::jsonb,
    true
),
-- CreArte
(
    'ac000000-0000-0000-0000-000000000009',
    'ac000000-0000-0000-0000-000000000008',
    'b2000000-0000-0000-0000-000000000002',
    'Grupo Manana',
    'GRP-MAN',
    'class',
    'Grupo de la manana - Modulo Pintura',
    'workshop',
    2026,
    '{"schedule": "morning"}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000011',
    'ac000000-0000-0000-0000-000000000010',
    'b2000000-0000-0000-0000-000000000002',
    'Grupo Tarde',
    'GRP-TAR',
    'class',
    'Grupo de la tarde - Modulo Escultura',
    'workshop',
    2026,
    '{"schedule": "afternoon"}'::jsonb,
    true
),
-- Academia
(
    'ac000000-0000-0000-0000-000000000014',
    'ac000000-0000-0000-0000-000000000013',
    'b3000000-0000-0000-0000-000000000003',
    'Class Monday',
    'CLS-MON',
    'class',
    'Monday class - Level A2',
    'language',
    2026,
    '{"day": "monday"}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000016',
    'ac000000-0000-0000-0000-000000000015',
    'b3000000-0000-0000-0000-000000000003',
    'Class Tuesday',
    'CLS-TUE',
    'class',
    'Tuesday class - Level B1',
    'language',
    2026,
    '{"day": "tuesday"}'::jsonb,
    true
)
ON CONFLICT (school_id, code, academic_year) DO UPDATE SET
    name           = EXCLUDED.name,
    parent_unit_id = EXCLUDED.parent_unit_id,
    description    = EXCLUDED.description,
    level          = EXCLUDED.level,
    metadata       = EXCLUDED.metadata,
    is_active      = EXCLUDED.is_active,
    updated_at     = now();

COMMIT;
