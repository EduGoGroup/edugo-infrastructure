-- Migration: 013_create_subjects
-- Description: Crea la tabla subjects para almacenar materias/asignaturas
-- Created: 2025-11-22
-- Project: edugo-infrastructure

CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índices para mejorar rendimiento
CREATE INDEX idx_subjects_name ON subjects(name);
CREATE INDEX idx_subjects_is_active ON subjects(is_active);
CREATE INDEX idx_subjects_created_at ON subjects(created_at DESC);
CREATE INDEX idx_subjects_metadata ON subjects USING GIN (metadata);

-- Comentarios para documentación
COMMENT ON TABLE subjects IS 'Materias o asignaturas del sistema educativo';
COMMENT ON COLUMN subjects.id IS 'Identificador único de la materia';
COMMENT ON COLUMN subjects.name IS 'Nombre de la materia';
COMMENT ON COLUMN subjects.description IS 'Descripción detallada de la materia';
COMMENT ON COLUMN subjects.metadata IS 'Metadata adicional en formato JSON';
COMMENT ON COLUMN subjects.is_active IS 'Indica si la materia está activa';
COMMENT ON COLUMN subjects.created_at IS 'Fecha de creación';
COMMENT ON COLUMN subjects.updated_at IS 'Fecha de última actualización';
