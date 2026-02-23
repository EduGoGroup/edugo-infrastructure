-- =============================================================================
-- EduGo Development Seeds — 007_materials.sql
-- =============================================================================
-- Crea 4 materiales educativos de prueba.
--
-- NOTA sobre el campo status (CHECK constraint en DB):
--   'uploaded'   → archivo subido, no procesado aún
--   'processing' → worker está procesando el material
--   'ready'      → procesado correctamente, assessment disponible
--   'failed'     → procesamiento fallido
--
-- Los materiales mat001, mat002, mat003 tienen status='ready' (ya procesados
-- y con assessment asociado en el siguiente seed).
-- El material mat004 tiene status='processing' (simulando un upload reciente
-- sin assessment generado aún).
--
-- Mapa de materiales:
--   mat001 → "Introducción a las Fracciones"    — math,    Clase 1-A, public
--   mat002 → "El Sistema Solar"                 — science, Clase 1-B, public
--   mat003 → "Historia de América Latina"        — history, Clase 10-A, privado
--   mat004 → "Álgebra Básica: Ecuaciones..."    — math,    Clase 1-A, privado, en proceso
-- =============================================================================

BEGIN;

INSERT INTO content.materials (
    id,
    school_id,
    uploaded_by_teacher_id,
    academic_unit_id,
    title,
    description,
    subject,
    grade,
    file_url,
    file_type,
    file_size_bytes,
    status,
    processing_started_at,
    processing_completed_at,
    is_public
) VALUES

-- -------------------------------------------------------------------------
-- Matemáticas — Clase 1-A — Escuela Primaria
-- -------------------------------------------------------------------------
(
    'aa100000-0000-0000-0000-000000000001',
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    '00000000-0000-0000-0000-000000000005',   -- teacher.math (María García)
    'ac000000-0000-0000-0000-000000000003',   -- Clase 1-A
    'Introducción a las Fracciones',
    'Material introductorio sobre fracciones simples, equivalentes y operaciones básicas para primer grado.',
    'Matemáticas',
    '1er Grado',
    's3://edugo-dev/materials/mat001.pdf',
    'application/pdf',
    2048000,
    'ready',
    '2024-09-10 10:00:00+00',
    '2024-09-10 10:05:32+00',
    true
),

-- -------------------------------------------------------------------------
-- Ciencias Naturales — Clase 1-B — Escuela Primaria
-- -------------------------------------------------------------------------
(
    'aa100000-0000-0000-0000-000000000002',
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    '00000000-0000-0000-0000-000000000006',   -- teacher.science (Juan Martínez)
    'ac000000-0000-0000-0000-000000000004',   -- Clase 1-B
    'El Sistema Solar',
    'Descripción de los planetas, el Sol y sus características principales, adaptada para primer grado.',
    'Ciencias Naturales',
    '1er Grado',
    's3://edugo-dev/materials/mat002.pdf',
    'application/pdf',
    3145728,
    'ready',
    '2024-09-12 11:00:00+00',
    '2024-09-12 11:04:18+00',
    true
),

-- -------------------------------------------------------------------------
-- Historia — Clase 10-A — Colegio Secundario
-- -------------------------------------------------------------------------
(
    'aa100000-0000-0000-0000-000000000003',
    'b2000000-0000-0000-0000-000000000002',   -- Colegio Secundario Demo
    '00000000-0000-0000-0000-000000000007',   -- teacher.history (Ana López)
    'ac000000-0000-0000-0000-000000000009',   -- Clase 10-A
    'Historia de América Latina',
    'Resumen de los principales procesos históricos de América Latina en los siglos XIX y XX.',
    'Historia',
    '10mo Grado',
    's3://edugo-dev/materials/mat003.pdf',
    'application/pdf',
    5242880,
    'ready',
    '2024-09-15 14:00:00+00',
    '2024-09-15 14:06:45+00',
    false
),

-- -------------------------------------------------------------------------
-- Matemáticas — Clase 1-A — en procesamiento (sin assessment aún)
-- -------------------------------------------------------------------------
(
    'aa100000-0000-0000-0000-000000000004',
    'b1000000-0000-0000-0000-000000000001',   -- Escuela Primaria Demo
    '00000000-0000-0000-0000-000000000005',   -- teacher.math (María García)
    'ac000000-0000-0000-0000-000000000003',   -- Clase 1-A
    'Álgebra Básica: Ecuaciones de Primer Grado',
    'Introducción a ecuaciones lineales simples con una incógnita, ejercicios prácticos y resolución paso a paso.',
    'Matemáticas',
    '1er Grado',
    's3://edugo-dev/materials/mat004.pdf',
    'application/pdf',
    1572864,
    'processing',
    NOW() - INTERVAL '5 minutes',
    NULL,
    false
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
