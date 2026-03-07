-- =============================================================================
-- EduGo Development Seeds v2 — 012_school_concepts.sql
-- =============================================================================
-- Copia concept_definitions a school_concepts para las 3 escuelas de dev.
-- Simula lo que hace SchoolService.CreateSchool en produccion.
--
-- Mapeo:
--   San Ignacio    → high_school       (c1000...002)
--   CreArte        → workshop          (c1000...005)
--   Academia       → language_academy  (c1000...003)
-- =============================================================================

BEGIN;

-- Colegio San Ignacio (high_school)
INSERT INTO academic.school_concepts (school_id, term_key, term_value, category)
SELECT 'b1000000-0000-0000-0000-000000000001', term_key, term_value, category
FROM academic.concept_definitions
WHERE concept_type_id = 'c1000000-0000-0000-0000-000000000002'
ON CONFLICT (school_id, term_key) DO NOTHING;

-- Taller CreArte (workshop)
INSERT INTO academic.school_concepts (school_id, term_key, term_value, category)
SELECT 'b2000000-0000-0000-0000-000000000002', term_key, term_value, category
FROM academic.concept_definitions
WHERE concept_type_id = 'c1000000-0000-0000-0000-000000000005'
ON CONFLICT (school_id, term_key) DO NOTHING;

-- Academia Global English (language_academy)
INSERT INTO academic.school_concepts (school_id, term_key, term_value, category)
SELECT 'b3000000-0000-0000-0000-000000000003', term_key, term_value, category
FROM academic.concept_definitions
WHERE concept_type_id = 'c1000000-0000-0000-0000-000000000003'
ON CONFLICT (school_id, term_key) DO NOTHING;

COMMIT;
