-- =============================================================================
-- EduGo Development Seeds — 010_guardian_relations.sql
-- =============================================================================
-- Crea 5 relaciones tutor-estudiante de prueba.
--
-- Mapa de relaciones:
--   gr001 → Roberto Gonzalez (u012) → Carlos Gonzalez (u008) — parent, primary, active
--   gr002 → Roberto Gonzalez (u012) → Sofia Rodriguez (u009) — guardian, secondary, active
--   gr003 → Patricia Torres  (u013) → Miguel Torres   (u010) — parent, primary, active
--   gr004 → Patricia Torres  (u013) → Laura Martinez  (u011) — guardian, secondary, active
--   gr005 → Roberto Gonzalez (u012) → Laura Martinez  (u011) — tutor, secondary, pending
-- =============================================================================

BEGIN;

INSERT INTO academic.guardian_relations (id, guardian_id, student_id, relationship_type, is_primary, is_active, status) VALUES
-- Roberto González (guardian u012) → Carlos González (student u008)
('ee000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000008', 'parent', true, true, 'active'),
-- Roberto González (guardian u012) → Sofía Rodríguez (student u009)
('ee000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000009', 'guardian', false, true, 'active'),
-- Patricia Torres (guardian u013) → Miguel Torres (student u010)
('ee000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000010', 'parent', true, true, 'active'),
-- Patricia Torres (guardian u013) → Laura Martínez (student u011)
('ee000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000011', 'guardian', false, true, 'active'),
-- Relación pendiente de prueba: Roberto González → Laura Martínez (pending)
('ee000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000011', 'tutor', false, true, 'pending')
ON CONFLICT (id) DO NOTHING;

COMMIT;
