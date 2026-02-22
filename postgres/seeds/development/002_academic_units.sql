-- =============================================================================
-- EduGo Development Seeds — 002_academic_units.sql
-- =============================================================================
-- Crea 9 unidades académicas en jerarquía correcta (padres antes que hijos).
--
-- Jerarquía Escuela Primaria (b1000...001):
--   au001 school  → "Escuela Primaria Demo"      (raíz, parent=NULL)
--   au002 grade   → "Primer Grado"               (parent=au001)
--   au003 class   → "Clase 1-A"                  (parent=au002)
--   au004 class   → "Clase 1-B"                  (parent=au002)
--   au005 grade   → "Segundo Grado"              (parent=au001)
--   au006 class   → "Clase 2-A"                  (parent=au005)
--
-- Jerarquía Colegio Secundario (b2000...002):
--   au007 school  → "Colegio Secundario Demo"    (raíz, parent=NULL)
--   au008 grade   → "Décimo Grado"               (parent=au007)
--   au009 class   → "Clase 10-A"                 (parent=au008)
--
-- NOTA sobre la constraint unique (school_id, code, academic_year):
--   Las unidades tipo 'school' usan academic_year=0 (valor por defecto).
--   Los grados y clases usan academic_year=2024.
-- =============================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- Nivel 0: Unidades raíz (type=school, parent=NULL)
-- -------------------------------------------------------------------------
INSERT INTO public.academic_units (
    id,
    parent_unit_id,
    school_id,
    name,
    code,
    type,
    description,
    level,
    academic_year,
    metadata,
    is_active
) VALUES
(
    'ac000000-0000-0000-0000-000000000001',
    NULL,
    'b1000000-0000-0000-0000-000000000001',
    'Escuela Primaria Demo',
    'EPD-ROOT',
    'school',
    'Unidad raíz de la Escuela Primaria Demo',
    'primary',
    0,
    '{"is_root": true}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000007',
    NULL,
    'b2000000-0000-0000-0000-000000000002',
    'Colegio Secundario Demo',
    'CSD-ROOT',
    'school',
    'Unidad raíz del Colegio Secundario Demo',
    'secondary',
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
-- Nivel 1: Grados (type=grade, parent=school)
-- -------------------------------------------------------------------------
INSERT INTO public.academic_units (
    id,
    parent_unit_id,
    school_id,
    name,
    code,
    type,
    description,
    level,
    academic_year,
    metadata,
    is_active
) VALUES
(
    'ac000000-0000-0000-0000-000000000002',
    'ac000000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',
    'Primer Grado',
    'GRADE-01',
    'grade',
    'Primer grado de educación primaria, año 2024',
    'primary',
    2024,
    '{"grade_number": 1}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000005',
    'ac000000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',
    'Segundo Grado',
    'GRADE-02',
    'grade',
    'Segundo grado de educación primaria, año 2024',
    'primary',
    2024,
    '{"grade_number": 2}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000008',
    'ac000000-0000-0000-0000-000000000007',
    'b2000000-0000-0000-0000-000000000002',
    'Décimo Grado',
    'GRADE-10',
    'grade',
    'Décimo grado de educación secundaria, año 2024',
    'secondary',
    2024,
    '{"grade_number": 10}'::jsonb,
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
-- Nivel 2: Clases (type=class, parent=grade)
-- -------------------------------------------------------------------------
INSERT INTO public.academic_units (
    id,
    parent_unit_id,
    school_id,
    name,
    code,
    type,
    description,
    level,
    academic_year,
    metadata,
    is_active
) VALUES
(
    'ac000000-0000-0000-0000-000000000003',
    'ac000000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000001',
    'Clase 1-A',
    '1A',
    'class',
    'Clase 1-A del Primer Grado, año 2024',
    'primary',
    2024,
    '{"section": "A", "grade_number": 1}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000004',
    'ac000000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000001',
    'Clase 1-B',
    '1B',
    'class',
    'Clase 1-B del Primer Grado, año 2024',
    'primary',
    2024,
    '{"section": "B", "grade_number": 1}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000006',
    'ac000000-0000-0000-0000-000000000005',
    'b1000000-0000-0000-0000-000000000001',
    'Clase 2-A',
    '2A',
    'class',
    'Clase 2-A del Segundo Grado, año 2024',
    'primary',
    2024,
    '{"section": "A", "grade_number": 2}'::jsonb,
    true
),
(
    'ac000000-0000-0000-0000-000000000009',
    'ac000000-0000-0000-0000-000000000008',
    'b2000000-0000-0000-0000-000000000002',
    'Clase 10-A',
    '10A',
    'class',
    'Clase 10-A del Décimo Grado, año 2024',
    'secondary',
    2024,
    '{"section": "A", "grade_number": 10}'::jsonb,
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
