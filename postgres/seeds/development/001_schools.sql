-- =============================================================================
-- EduGo Development Seeds v2 — 001_schools.sql
-- =============================================================================
-- 3 instituciones con tipos de concepto diferentes para probar terminologia:
--   b1000...001 → Colegio San Ignacio    (high_school)   — premium
--   b2000...002 → Taller CreArte         (workshop)      — basic
--   b3000...003 → Academia Global English (language_academy) — basic
--
-- Referencia concept_type_id (produccion 008_concept_types.sql):
--   c1000...002 → high_school
--   c1000...005 → workshop
--   c1000...003 → language_academy
-- =============================================================================

BEGIN;

INSERT INTO academic.schools (
    id,
    name,
    code,
    address,
    city,
    country,
    phone,
    email,
    concept_type_id,
    metadata,
    is_active,
    subscription_tier,
    max_teachers,
    max_students
) VALUES
-- Colegio San Ignacio (high_school)
(
    'b1000000-0000-0000-0000-000000000001',
    'Colegio San Ignacio',
    'SCH_SI_001',
    'Av. Libertador 1500',
    'Santiago',
    'Chile',
    '+56 2 2345 6789',
    'contacto@sanignacio.edugo.test',
    'c1000000-0000-0000-0000-000000000002',
    '{"level": "secondary", "demo": true, "founded_year": 2018}'::jsonb,
    true,
    'premium',
    20,
    300
),
-- Taller CreArte (workshop)
(
    'b2000000-0000-0000-0000-000000000002',
    'Taller CreArte',
    'SCH_CA_001',
    'Calle Artistas 234',
    'Valparaiso',
    'Chile',
    '+56 32 2345 678',
    'contacto@crearte.edugo.test',
    'c1000000-0000-0000-0000-000000000005',
    '{"level": "workshop", "demo": true, "founded_year": 2021}'::jsonb,
    true,
    'basic',
    10,
    100
),
-- Academia Global English (language_academy)
(
    'b3000000-0000-0000-0000-000000000003',
    'Academia Global English',
    'SCH_GE_001',
    'Paseo Internacional 89',
    'Santiago',
    'Chile',
    '+56 2 9876 5432',
    'contacto@globalenglish.edugo.test',
    'c1000000-0000-0000-0000-000000000003',
    '{"level": "language_academy", "demo": true, "founded_year": 2020}'::jsonb,
    true,
    'basic',
    10,
    150
)
ON CONFLICT (id) DO UPDATE SET
    name              = EXCLUDED.name,
    code              = EXCLUDED.code,
    address           = EXCLUDED.address,
    city              = EXCLUDED.city,
    country           = EXCLUDED.country,
    phone             = EXCLUDED.phone,
    email             = EXCLUDED.email,
    concept_type_id   = EXCLUDED.concept_type_id,
    metadata          = EXCLUDED.metadata,
    is_active         = EXCLUDED.is_active,
    subscription_tier = EXCLUDED.subscription_tier,
    max_teachers      = EXCLUDED.max_teachers,
    max_students      = EXCLUDED.max_students,
    updated_at        = now();

COMMIT;
