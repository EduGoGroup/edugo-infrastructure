-- Migration: 014_create_units
-- Description: Crea la tabla units para almacenar unidades organizacionales (departamentos, grados, grupos, etc.)
-- Created: 2025-11-22
-- Project: edugo-infrastructure

CREATE TABLE IF NOT EXISTS units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    parent_unit_id UUID REFERENCES units(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL CHECK (length(name) >= 2),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    CHECK (id != parent_unit_id) -- No puede ser su propio padre
);

-- Índices para mejorar rendimiento
CREATE INDEX idx_units_school_id ON units(school_id);
CREATE INDEX idx_units_parent_unit_id ON units(parent_unit_id);
CREATE INDEX idx_units_name ON units(name);
CREATE INDEX idx_units_is_active ON units(is_active);
CREATE INDEX idx_units_created_at ON units(created_at DESC);

-- Índice para consultas jerárquicas
CREATE INDEX idx_units_hierarchy ON units(school_id, parent_unit_id, is_active);

-- Comentarios para documentación
COMMENT ON TABLE units IS 'Unidades organizacionales jerárquicas (departamentos, grados, grupos, etc.)';
COMMENT ON COLUMN units.id IS 'Identificador único de la unidad';
COMMENT ON COLUMN units.school_id IS 'ID de la escuela a la que pertenece';
COMMENT ON COLUMN units.parent_unit_id IS 'ID de la unidad padre (NULL si es raíz)';
COMMENT ON COLUMN units.name IS 'Nombre de la unidad (mínimo 2 caracteres)';
COMMENT ON COLUMN units.description IS 'Descripción de la unidad';
COMMENT ON COLUMN units.is_active IS 'Indica si la unidad está activa';
COMMENT ON COLUMN units.created_at IS 'Fecha de creación';
COMMENT ON COLUMN units.updated_at IS 'Fecha de última actualización';
