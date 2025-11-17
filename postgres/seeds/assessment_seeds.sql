-- Seeds: Sistema de Evaluaciones
-- Description: Datos de prueba para desarrollo y testing cross-ecosystem
-- Dependencies: Migraciones 006-011, tabla materials y users existentes
-- Date: 2025-11-17

BEGIN;

-- Verificar prerequisitos: tabla assessment
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'assessment') THEN
        RAISE EXCEPTION 'Tabla assessment no existe. Ejecutar migraciones primero.';
    END IF;
END $$;

-- Verificar prerequisitos: materials suficientes
DO $$
DECLARE
    v_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_count FROM materials;
    IF v_count < 3 THEN
        RAISE EXCEPTION 'Se requieren al menos 3 materials en la base de datos. Actual: %', v_count;
    END IF;
END $$;

-- Protección contra producción
DO $$
BEGIN
    IF current_database() = 'edugo_prod' THEN
        RAISE EXCEPTION 'SEEDS PROHIBIDOS EN PRODUCCIÓN';
    END IF;
END $$;

TRUNCATE TABLE assessment_attempt_answer CASCADE;
TRUNCATE TABLE assessment_attempt CASCADE;
TRUNCATE TABLE assessment CASCADE;

-- Seed 1: Assessment de Programación
INSERT INTO assessment (
    id, material_id, mongo_document_id, title,
    total_questions, questions_count, pass_threshold,
    max_attempts, time_limit_minutes, status,
    created_at, updated_at
) VALUES (
    gen_random_uuid(),
    (SELECT id FROM materials LIMIT 1),
    '507f1f77bcf86cd799439011',
    'Quiz: Fundamentos de Programación',
    10, 10, 70, NULL, 30, 'published',
    NOW() - INTERVAL '30 days',
    NOW() - INTERVAL '30 days'
) ON CONFLICT (mongo_document_id) DO NOTHING;

-- Seed 2: Assessment de Algoritmos
INSERT INTO assessment (
    id, material_id, mongo_document_id, title,
    total_questions, questions_count, pass_threshold,
    max_attempts, time_limit_minutes, status,
    created_at, updated_at
) VALUES (
    gen_random_uuid(),
    (SELECT id FROM materials LIMIT 1 OFFSET 1),
    '507f1f77bcf86cd799439012',
    'Quiz: Algoritmos de Ordenamiento',
    8, 8, 75, 3, 20, 'published',
    NOW() - INTERVAL '15 days',
    NOW() - INTERVAL '15 days'
) ON CONFLICT (mongo_document_id) DO NOTHING;

-- Seed 3: Assessment en draft
INSERT INTO assessment (
    id, material_id, mongo_document_id, title,
    total_questions, questions_count, pass_threshold,
    max_attempts, time_limit_minutes, status,
    created_at, updated_at
) VALUES (
    gen_random_uuid(),
    (SELECT id FROM materials LIMIT 1 OFFSET 2),
    '507f1f77bcf86cd799439013',
    'Quiz: Estructuras de Datos (Borrador)',
    12, 12, 80, NULL, 40, 'draft',
    NOW() - INTERVAL '7 days',
    NOW() - INTERVAL '7 days'
) ON CONFLICT (mongo_document_id) DO NOTHING;

COMMIT;

-- Verificación: comprobar que los IDs específicos existen
DO $$
DECLARE
    v_missing TEXT := '';
BEGIN
    IF NOT EXISTS (SELECT 1 FROM assessment WHERE mongo_document_id = '507f1f77bcf86cd799439011') THEN
        v_missing := v_missing || '507f1f77bcf86cd799439011 ';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM assessment WHERE mongo_document_id = '507f1f77bcf86cd799439012') THEN
        v_missing := v_missing || '507f1f77bcf86cd799439012 ';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM assessment WHERE mongo_document_id = '507f1f77bcf86cd799439013') THEN
        v_missing := v_missing || '507f1f77bcf86cd799439013 ';
    END IF;
    
    IF v_missing <> '' THEN
        RAISE EXCEPTION 'Faltan assessments con mongo_document_id: %', v_missing;
    ELSE
        RAISE NOTICE 'Todos los assessments esperados están presentes.';
    END IF;
END $$;
