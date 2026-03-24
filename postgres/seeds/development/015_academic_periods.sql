-- =============================================================================
-- EduGo Development Seeds v2 — 015_academic_periods.sql
-- =============================================================================
-- Periodos academicos de ejemplo para las 3 escuelas.
-- =============================================================================

BEGIN;

INSERT INTO academic.academic_periods (id, school_id, name, code, type, start_date, end_date, is_active, academic_year, sort_order) VALUES
-- San Ignacio - 2 semestres
('ff000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', 'Primer Semestre 2026', 'S1-2026', 'semester', '2026-03-01', '2026-07-15', true, 2026, 1),
('ff000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001', 'Segundo Semestre 2026', 'S2-2026', 'semester', '2026-08-01', '2026-12-15', false, 2026, 2),
-- CreArte - trimestres
('ff000000-0000-0000-0000-000000000003', 'b2000000-0000-0000-0000-000000000002', 'Primer Trimestre 2026', 'T1-2026', 'trimester', '2026-03-01', '2026-05-31', true, 2026, 1),
('ff000000-0000-0000-0000-000000000004', 'b2000000-0000-0000-0000-000000000002', 'Segundo Trimestre 2026', 'T2-2026', 'trimester', '2026-06-01', '2026-08-31', false, 2026, 2),
-- Academia - bimestres
('ff000000-0000-0000-0000-000000000005', 'b3000000-0000-0000-0000-000000000003', 'Bimestre 1', 'B1-2026', 'bimester', '2026-03-01', '2026-04-30', true, 2026, 1)
ON CONFLICT (id) DO NOTHING;

COMMIT;
