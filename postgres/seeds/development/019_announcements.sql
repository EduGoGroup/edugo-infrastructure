-- =============================================================================
-- EduGo Development Seeds v2 — 019_announcements.sql
-- =============================================================================
-- Anuncios de ejemplo para San Ignacio.
-- 2 anuncios: uno a nivel escuela, otro a nivel unidad 5to A.
-- =============================================================================

BEGIN;

INSERT INTO academic.announcements (id, school_id, academic_unit_id, author_id, title, body, scope, target_roles, is_pinned, published_at, expires_at) VALUES
-- Anuncio escolar general (por Carmen, school_admin)
('a3000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', NULL, '00000000-0000-0000-0000-000000000002', 'Reunion de Apoderados Marzo 2026', 'Se convoca a reunion general de apoderados el viernes 28 de marzo a las 18:00 hrs en el auditorio principal. Favor confirmar asistencia.', 'school', '{guardian,teacher}', true, '2026-03-20 10:00:00+00', '2026-03-28 23:59:59+00'),
-- Anuncio de unidad 5to A (por Maria, teacher)
('a3000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000005', 'Prueba de Matematicas - 5to A', 'La prueba de Matematicas del primer semestre se realizara el lunes 31 de marzo. Estudiar capitulos 1 al 4.', 'unit', '{student,guardian}', false, '2026-03-21 14:00:00+00', '2026-03-31 23:59:59+00')
ON CONFLICT (id) DO NOTHING;

COMMIT;
