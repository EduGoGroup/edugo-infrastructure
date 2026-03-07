-- =============================================================================
-- EduGo Development Seeds v2 — 007_materials.sql
-- =============================================================================
-- 5 materiales educativos vinculados a las 3 escuelas.
--
-- Mapa:
--   mat001 → Introduccion a las Fracciones   → San Ignacio → Maria (U-05) → 5to A → ready
--   mat002 → El Sistema Solar                → San Ignacio → Pedro (U-06) → 5to A → ready
--   mat003 → Historia de Chile: Independencia→ San Ignacio → Pedro (U-06) → 6to A → ready
--   mat004 → Teoria del Color                → CreArte     → Ana (U-07)   → Grp Manana → ready
--   mat005 → English Grammar Basics          → Academia    → Maria (U-05) → Class Monday → ready
-- =============================================================================

BEGIN;

INSERT INTO content.materials (
    id, school_id, uploaded_by_teacher_id, academic_unit_id,
    title, description, subject, grade,
    file_url, file_type, file_size_bytes,
    status, processing_started_at, processing_completed_at, is_public
) VALUES
(
    'aa100000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000005',
    'ac000000-0000-0000-0000-000000000003',
    'Introduccion a las Fracciones',
    'Material introductorio sobre fracciones simples, equivalentes y operaciones basicas.',
    'Matematicas', '5to Basico',
    's3://edugo-dev/materials/mat001.pdf', 'application/pdf', 2048000,
    'ready', '2026-02-10 10:00:00+00', '2026-02-10 10:05:32+00', true
),
(
    'aa100000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000006',
    'ac000000-0000-0000-0000-000000000003',
    'El Sistema Solar',
    'Descripcion de los planetas, el Sol y sus caracteristicas principales.',
    'Ciencias Naturales', '5to Basico',
    's3://edugo-dev/materials/mat002.pdf', 'application/pdf', 3145728,
    'ready', '2026-02-12 11:00:00+00', '2026-02-12 11:04:18+00', true
),
(
    'aa100000-0000-0000-0000-000000000003',
    'b1000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000006',
    'ac000000-0000-0000-0000-000000000006',
    'Historia de Chile: Independencia',
    'Resumen de los principales procesos de la independencia de Chile.',
    'Historia', '6to Basico',
    's3://edugo-dev/materials/mat003.pdf', 'application/pdf', 5242880,
    'ready', '2026-02-15 14:00:00+00', '2026-02-15 14:06:45+00', false
),
(
    'aa100000-0000-0000-0000-000000000004',
    'b2000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000007',
    'ac000000-0000-0000-0000-000000000009',
    'Teoria del Color',
    'Fundamentos de la teoria del color: colores primarios, secundarios, complementarios.',
    'Pintura', 'Modulo Pintura',
    's3://edugo-dev/materials/mat004.pdf', 'application/pdf', 1800000,
    'ready', '2026-02-14 09:00:00+00', '2026-02-14 09:03:22+00', true
),
(
    'aa100000-0000-0000-0000-000000000005',
    'b3000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000005',
    'ac000000-0000-0000-0000-000000000014',
    'English Grammar Basics',
    'Introduction to basic English grammar: articles, pronouns, simple tenses.',
    'English', 'Level A2',
    's3://edugo-dev/materials/mat005.pdf', 'application/pdf', 1500000,
    'ready', '2026-02-16 10:00:00+00', '2026-02-16 10:04:10+00', true
)
ON CONFLICT (id) DO UPDATE SET
    title                   = EXCLUDED.title,
    description             = EXCLUDED.description,
    subject                 = EXCLUDED.subject,
    grade                   = EXCLUDED.grade,
    file_url                = EXCLUDED.file_url,
    file_type               = EXCLUDED.file_type,
    file_size_bytes         = EXCLUDED.file_size_bytes,
    status                  = EXCLUDED.status,
    processing_started_at   = EXCLUDED.processing_started_at,
    processing_completed_at = EXCLUDED.processing_completed_at,
    is_public               = EXCLUDED.is_public,
    updated_at              = now();

COMMIT;
