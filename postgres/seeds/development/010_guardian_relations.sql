-- =============================================================================
-- EduGo Development Seeds v2 — 010_guardian_relations.sql
-- =============================================================================
-- 3 relaciones tutor-estudiante.
--
-- Mapa:
--   gr001 → Ricardo Mendoza (U-13) → Carlos Mendoza (U-08) — parent, primary
--   gr002 → Patricia Herrera (U-14) → Sofia Herrera (U-09) — parent, primary
--   gr003 → Patricia Herrera (U-14) → Diego Vargas (U-10)  — guardian, secondary
-- =============================================================================

BEGIN;

INSERT INTO academic.guardian_relations (id, guardian_id, student_id, relationship_type, is_primary, is_active, status) VALUES
('ee000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000008', 'parent', true, true, 'active'),
('ee000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000009', 'parent', true, true, 'active'),
('ee000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000010', 'guardian', false, true, 'active')
ON CONFLICT (id) DO NOTHING;

COMMIT;
