-- =============================================================================
-- EduGo Development Seeds v2 — 010_guardian_relations.sql
-- =============================================================================
-- 5 relaciones tutor-estudiante.
--
-- Mapa:
--   gr001 → Ricardo Mendoza (U-13)  → Carlos Mendoza (U-08)  — parent, primary, active
--   gr002 → Patricia Herrera (U-14) → Sofia Herrera (U-09)   — parent, primary, active
--   gr003 → Patricia Herrera (U-14) → Diego Vargas (U-10)    — guardian, secondary, active
--   gr004 → Fernando Ruiz (U-20)    → Carlos Mendoza (U-08)  — guardian, secondary, pending
--   gr005 → Patricia Herrera (U-14) → Valentina Rojas (U-11) — guardian, secondary, pending
-- =============================================================================

BEGIN;

INSERT INTO academic.guardian_relations (id, guardian_id, student_id, relationship_type, is_primary, is_active, status) VALUES
('ee000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000008', 'parent', true, true, 'active'),
('ee000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000009', 'parent', true, true, 'active'),
('ee000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000010', 'guardian', false, true, 'active'),
-- gr004: Fernando Ruiz (U-20) → Carlos Mendoza (U-08) — tio, PENDING
('ee000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000020', '00000000-0000-0000-0000-000000000008', 'guardian', false, true, 'pending'),
-- gr005: Patricia Herrera (U-14) → Valentina Rojas (U-11) — tutora, PENDING
('ee000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000011', 'guardian', false, true, 'pending')
ON CONFLICT (id) DO NOTHING;

COMMIT;
