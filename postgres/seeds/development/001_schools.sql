-- =============================================================================
-- EduGo Development Seeds — 001_schools.sql
-- =============================================================================
-- Crea 3 escuelas de demo representando distintos niveles educativos:
--   b1000...001 → Escuela Primaria Demo   (subscription: premium)
--   b2000...002 → Colegio Secundario Demo (subscription: basic)
--   b3000...003 → Instituto Técnico Demo  (subscription: free)
--
-- IDs fijos para que los seeds posteriores puedan referenciarlos sin
-- consultar la DB.
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
    metadata,
    is_active,
    subscription_tier,
    max_teachers,
    max_students
) VALUES
-- Escuela Primaria Demo
(
    'b1000000-0000-0000-0000-000000000001',
    'Escuela Primaria Demo',
    'SCH_PRI_001',
    'Av. Educación 1234',
    'Santiago',
    'Chile',
    '+56 2 2345 6789',
    'contacto@primaria.demo.edugo.test',
    '{"level": "primary", "demo": true, "founded_year": 2020}'::jsonb,
    true,
    'premium',
    20,
    300
),
-- Colegio Secundario Demo
(
    'b2000000-0000-0000-0000-000000000002',
    'Colegio Secundario Demo',
    'SCH_SEC_001',
    'Calle Ciencia 567',
    'Valparaíso',
    'Chile',
    '+56 32 2345 678',
    'contacto@secundario.demo.edugo.test',
    '{"level": "secondary", "demo": true, "founded_year": 2021}'::jsonb,
    true,
    'basic',
    15,
    200
),
-- Instituto Técnico Demo
(
    'b3000000-0000-0000-0000-000000000003',
    'Instituto Técnico Demo',
    'SCH_TEC_001',
    'Pasaje Técnico 89',
    'Concepción',
    'Chile',
    '+56 41 2345 678',
    'contacto@tecnico.demo.edugo.test',
    '{"level": "technical", "demo": true, "founded_year": 2022}'::jsonb,
    true,
    'free',
    10,
    100
)
ON CONFLICT (id) DO UPDATE SET
    name              = EXCLUDED.name,
    code              = EXCLUDED.code,
    address           = EXCLUDED.address,
    city              = EXCLUDED.city,
    country           = EXCLUDED.country,
    phone             = EXCLUDED.phone,
    email             = EXCLUDED.email,
    metadata          = EXCLUDED.metadata,
    is_active         = EXCLUDED.is_active,
    subscription_tier = EXCLUDED.subscription_tier,
    max_teachers      = EXCLUDED.max_teachers,
    max_students      = EXCLUDED.max_students,
    updated_at        = now();

COMMIT;
