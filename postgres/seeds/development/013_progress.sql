-- =============================================================================
-- EduGo Development Seeds v2 — 013_progress.sql
-- =============================================================================
-- 8 registros de progreso de estudiantes en materiales educativos.
--
-- Tabla: content.progress (PK compuesta: material_id, user_id)
-- Columnas: material_id, user_id, progress_percentage, last_position, completed_at
--
-- Mapa:
--   pr001 → Carlos (U-08)    → Fracciones (mat001)       → 100% (completado)
--   pr002 → Carlos (U-08)    → Sistema Solar (mat002)     →  65%
--   pr003 → Carlos (U-08)    → Teoria Color (mat004)      →  30%
--   pr004 → Sofia (U-09)     → Fracciones (mat001)        →  80%
--   pr005 → Sofia (U-09)     → Sistema Solar (mat002)     →  45%
--   pr006 → Diego (U-10)     → Sistema Solar (mat002)     →  90%
--   pr007 → Valentina (U-11) → English Grammar (mat005)   →  70%
--   pr008 → Mateo (U-12)     → Teoria Color (mat004)      →  55%
-- =============================================================================

BEGIN;

INSERT INTO content.progress (
    material_id,
    user_id,
    progress_percentage,
    last_position,
    completed_at
) VALUES

-- pr001: Carlos → Fracciones → 100% complete
(
    'aa100000-0000-0000-0000-000000000001',   -- Fracciones
    '00000000-0000-0000-0000-000000000008',   -- Carlos Mendoza
    100.00,
    '{"page": 24, "section": "ejercicios-finales"}'::jsonb,
    '2026-03-10 15:30:00+00'
),

-- pr002: Carlos → Sistema Solar → 65%
(
    'aa100000-0000-0000-0000-000000000002',   -- Sistema Solar
    '00000000-0000-0000-0000-000000000008',   -- Carlos Mendoza
    65.00,
    '{"page": 12, "section": "planetas-exteriores"}'::jsonb,
    NULL
),

-- pr003: Carlos → Teoria del Color → 30%
(
    'aa100000-0000-0000-0000-000000000004',   -- Teoria del Color
    '00000000-0000-0000-0000-000000000008',   -- Carlos Mendoza
    30.00,
    '{"page": 5, "section": "colores-primarios"}'::jsonb,
    NULL
),

-- pr004: Sofia → Fracciones → 80%
(
    'aa100000-0000-0000-0000-000000000001',   -- Fracciones
    '00000000-0000-0000-0000-000000000009',   -- Sofia Herrera
    80.00,
    '{"page": 19, "section": "fracciones-equivalentes"}'::jsonb,
    NULL
),

-- pr005: Sofia → Sistema Solar → 45%
(
    'aa100000-0000-0000-0000-000000000002',   -- Sistema Solar
    '00000000-0000-0000-0000-000000000009',   -- Sofia Herrera
    45.00,
    '{"page": 8, "section": "planetas-interiores"}'::jsonb,
    NULL
),

-- pr006: Diego → Sistema Solar → 90%
(
    'aa100000-0000-0000-0000-000000000002',   -- Sistema Solar
    '00000000-0000-0000-0000-000000000010',   -- Diego Vargas
    90.00,
    '{"page": 20, "section": "resumen"}'::jsonb,
    NULL
),

-- pr007: Valentina → English Grammar → 70%
(
    'aa100000-0000-0000-0000-000000000005',   -- English Grammar Basics
    '00000000-0000-0000-0000-000000000011',   -- Valentina Rojas
    70.00,
    '{"page": 14, "section": "simple-tenses"}'::jsonb,
    NULL
),

-- pr008: Mateo → Teoria del Color → 55%
(
    'aa100000-0000-0000-0000-000000000004',   -- Teoria del Color
    '00000000-0000-0000-0000-000000000012',   -- Mateo Fuentes
    55.00,
    '{"page": 9, "section": "colores-secundarios"}'::jsonb,
    NULL
)

ON CONFLICT (material_id, user_id) DO UPDATE SET
    progress_percentage = EXCLUDED.progress_percentage,
    last_position       = EXCLUDED.last_position,
    completed_at        = EXCLUDED.completed_at,
    updated_at          = now();

COMMIT;
