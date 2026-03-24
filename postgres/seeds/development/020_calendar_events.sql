-- =============================================================================
-- EduGo Development Seeds v2 — 020_calendar_events.sql
-- =============================================================================
-- Eventos de calendario de ejemplo para San Ignacio.
-- 2 eventos: un feriado y un examen.
-- =============================================================================

BEGIN;

DO $$ BEGIN
IF to_regclass('academic.calendar_events') IS NOT NULL THEN

INSERT INTO academic.calendar_events (id, school_id, title, description, event_type, start_date, end_date, is_all_day, created_by) VALUES
-- Feriado: Semana Santa
('a4000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', 'Semana Santa - Sin Clases', 'Receso por Semana Santa. Se retoman clases el lunes 6 de abril.', 'holiday', '2026-04-02', '2026-04-05', true, '00000000-0000-0000-0000-000000000002'),
-- Examen: Evaluaciones de Matematicas
('a4000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001', 'Examenes Primer Semestre - Matematicas', 'Examenes de matematicas para todos los cursos de 5to y 6to.', 'exam', '2026-03-31', '2026-03-31', true, '00000000-0000-0000-0000-000000000002')
ON CONFLICT (id) DO NOTHING;

END IF;
END $$;

COMMIT;
