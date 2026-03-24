-- =============================================================================
-- EduGo Development Seeds v2 — 014_membership_subjects.sql
-- =============================================================================
-- Vincula memberships de tipo teacher/assistant_teacher con sus materias.
-- Reemplaza los datos que antes estaban en memberships.metadata.subjects.
--
-- Mapa (membership → subject):
--   m008 (Maria, teacher 5A)          → sub001 (MAT-5A)
--   m009 (Maria, teacher Academia)    → sub007 (ENG-A2)
--   m010 (Pedro, teacher 5B)          → sub003 (MAT-5B)
--   m010 (Pedro, teacher 5B)          → sub008 (SCI-5B)
--   m011 (Pedro, teacher 6A)          → sub004 (HIS-6A)
--   m012 (Ana, teacher CreArte GM)    → sub005 (PINT-GM)
--   m013 (Ana, teacher CreArte GT)    → sub006 (ESCL-GT)
--   m024 (Andres, assistant 5A)       → sub001 (MAT-5A)
-- =============================================================================

BEGIN;

INSERT INTO academic.membership_subjects (membership_id, subject_id) VALUES
-- Maria Martinez → Matematicas en 5to A
('bb000000-0000-0000-0000-000000000008', 'dd000000-0000-0000-0000-000000000001'),
-- Maria Martinez → English Basics A2 en Academia
('bb000000-0000-0000-0000-000000000009', 'dd000000-0000-0000-0000-000000000007'),
-- Pedro Gonzalez → Matematicas en 5to B
('bb000000-0000-0000-0000-000000000010', 'dd000000-0000-0000-0000-000000000003'),
-- Pedro Gonzalez → Ciencias Naturales en 5to B
('bb000000-0000-0000-0000-000000000010', 'dd000000-0000-0000-0000-000000000008'),
-- Pedro Gonzalez → Historia en 6to A
('bb000000-0000-0000-0000-000000000011', 'dd000000-0000-0000-0000-000000000004'),
-- Ana Ruiz → Tecnicas de Pintura en Grupo Manana
('bb000000-0000-0000-0000-000000000012', 'dd000000-0000-0000-0000-000000000005'),
-- Ana Ruiz → Fundamentos de Escultura en Grupo Tarde
('bb000000-0000-0000-0000-000000000013', 'dd000000-0000-0000-0000-000000000006'),
-- Andres Gomez (assistant_teacher) → Matematicas en 5to A
('bb000000-0000-0000-0000-000000000024', 'dd000000-0000-0000-0000-000000000001')
ON CONFLICT (membership_id, subject_id) DO NOTHING;

COMMIT;
