-- =============================================================================
-- EduGo Development Seeds — 012_school_concepts.sql
-- =============================================================================
-- Copia concept_definitions a school_concepts para las escuelas de dev.
-- Simula lo que hace SchoolService.CreateSchool en produccion.
-- =============================================================================

BEGIN;

-- Escuela Primaria Demo (primary_school)
INSERT INTO academic.school_concepts (school_id, term_key, term_value, category)
SELECT 'b1000000-0000-0000-0000-000000000001', term_key, term_value, category
FROM academic.concept_definitions
WHERE concept_type_id = 'c1000000-0000-0000-0000-000000000001'
ON CONFLICT (school_id, term_key) DO NOTHING;

-- Colegio Secundario Demo (high_school)
INSERT INTO academic.school_concepts (school_id, term_key, term_value, category)
SELECT 'b2000000-0000-0000-0000-000000000002', term_key, term_value, category
FROM academic.concept_definitions
WHERE concept_type_id = 'c1000000-0000-0000-0000-000000000002'
ON CONFLICT (school_id, term_key) DO NOTHING;

-- Instituto Tecnico Demo (technical_school)
INSERT INTO academic.school_concepts (school_id, term_key, term_value, category)
SELECT 'b3000000-0000-0000-0000-000000000003', term_key, term_value, category
FROM academic.concept_definitions
WHERE concept_type_id = 'c1000000-0000-0000-0000-000000000004'
ON CONFLICT (school_id, term_key) DO NOTHING;

COMMIT;
