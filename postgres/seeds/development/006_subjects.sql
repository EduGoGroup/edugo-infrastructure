-- =============================================================================
-- EduGo Development Seeds v2 — 006_subjects.sql
-- =============================================================================
-- 8 materias/talleres/courses vinculados a las 3 escuelas.
--
-- Mapa:
--   sub001 → Matematicas            → San Ignacio → 5to A (au003)
--   sub002 → Ciencias Naturales     → San Ignacio → 5to A (au003)
--   sub003 → Matematicas            → San Ignacio → 5to B (au004)
--   sub004 → Historia               → San Ignacio → 6to A (au006)
--   sub005 → Tecnicas de Pintura    → CreArte     → Grupo Manana (au009)
--   sub006 → Fundamentos Escultura  → CreArte     → Grupo Tarde (au011)
--   sub007 → English Basics A2      → Academia    → Class Monday (au014)
--   sub008 → Ciencias Naturales     → San Ignacio → 5to B (au004)
-- =============================================================================

BEGIN;

INSERT INTO academic.subjects (id, school_id, academic_unit_id, name, code, description, is_active) VALUES
-- San Ignacio
('dd000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000003', 'Matematicas', 'MAT-5A', 'Matematicas para 5to A', true),
('dd000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000003', 'Ciencias Naturales', 'SCI-5A', 'Ciencias Naturales para 5to A', true),
('dd000000-0000-0000-0000-000000000003', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000004', 'Matematicas', 'MAT-5B', 'Matematicas para 5to B', true),
('dd000000-0000-0000-0000-000000000004', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000006', 'Historia', 'HIS-6A', 'Historia de Chile para 6to A', true),
-- CreArte
('dd000000-0000-0000-0000-000000000005', 'b2000000-0000-0000-0000-000000000002', 'ac000000-0000-0000-0000-000000000009', 'Tecnicas de Pintura', 'PINT-GM', 'Taller de tecnicas de pintura', true),
('dd000000-0000-0000-0000-000000000006', 'b2000000-0000-0000-0000-000000000002', 'ac000000-0000-0000-0000-000000000011', 'Fundamentos de Escultura', 'ESCL-GT', 'Taller de fundamentos de escultura', true),
-- Academia
('dd000000-0000-0000-0000-000000000007', 'b3000000-0000-0000-0000-000000000003', 'ac000000-0000-0000-0000-000000000014', 'English Basics A2', 'ENG-A2', 'English course for level A2', true),
-- San Ignacio (5to B)
('dd000000-0000-0000-0000-000000000008', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000004', 'Ciencias Naturales', 'SCI-5B', 'Ciencias Naturales para 5to B', true)
ON CONFLICT (id) DO NOTHING;

COMMIT;
