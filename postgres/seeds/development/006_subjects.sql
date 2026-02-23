-- =============================================================================
-- EduGo Development Seeds — 006_subjects.sql
-- =============================================================================
-- Crea 5 materias de prueba vinculadas a las escuelas demo.
--
-- Mapa de subjects:
--   sub001 → Matematicas        — Primaria, Clase 1-A (au004)
--   sub002 → Ciencias Naturales — Primaria, Clase 1-A (au004)
--   sub003 → Matematicas        — Primaria, Clase 1-B (au005)
--   sub004 → Historia           — Secundario, Decimo Grado (au007)
--   sub005 → Fisica             — Secundario, Decimo Grado (au007)
-- =============================================================================

BEGIN;

INSERT INTO academic.subjects (id, school_id, academic_unit_id, name, code, description, is_active) VALUES
-- Escuela Primaria (b1000000-...-001)
('dd000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000004', 'Matematicas', 'MAT-1A', 'Matematicas basicas para primer grado', true),
('dd000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000004', 'Ciencias Naturales', 'SCI-1A', 'Ciencias naturales para primer grado', true),
('dd000000-0000-0000-0000-000000000003', 'b1000000-0000-0000-0000-000000000001', 'ac000000-0000-0000-0000-000000000005', 'Matematicas', 'MAT-2A', 'Matematicas para segundo grado', true),
-- Escuela Secundaria (b2000000-...-002)
('dd000000-0000-0000-0000-000000000004', 'b2000000-0000-0000-0000-000000000002', 'ac000000-0000-0000-0000-000000000007', 'Historia', 'HIS-1', 'Historia universal', true),
('dd000000-0000-0000-0000-000000000005', 'b2000000-0000-0000-0000-000000000002', 'ac000000-0000-0000-0000-000000000007', 'Fisica', 'FIS-1', 'Fisica basica', true)
ON CONFLICT (id) DO NOTHING;

COMMIT;
