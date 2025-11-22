-- Migration: 012_create_material_versions
-- Description: Crea la tabla material_versions para almacenar versiones históricas de materiales
-- Created: 2025-11-22
-- Project: edugo-infrastructure

CREATE TABLE IF NOT EXISTS material_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL CHECK (version_number > 0),
    title VARCHAR(255) NOT NULL,
    content_url TEXT NOT NULL,
    changed_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    UNIQUE(material_id, version_number)
);

-- Índices para mejorar rendimiento
CREATE INDEX idx_material_versions_material_id ON material_versions(material_id);
CREATE INDEX idx_material_versions_version_number ON material_versions(material_id, version_number DESC);
CREATE INDEX idx_material_versions_changed_by ON material_versions(changed_by);
CREATE INDEX idx_material_versions_created_at ON material_versions(created_at DESC);

-- Comentarios para documentación
COMMENT ON TABLE material_versions IS 'Historial de versiones de materiales educativos';
COMMENT ON COLUMN material_versions.id IS 'Identificador único de la versión';
COMMENT ON COLUMN material_versions.material_id IS 'ID del material al que pertenece esta versión';
COMMENT ON COLUMN material_versions.version_number IS 'Número de versión (1, 2, 3...)';
COMMENT ON COLUMN material_versions.title IS 'Título de esta versión del material';
COMMENT ON COLUMN material_versions.content_url IS 'URL del contenido de esta versión';
COMMENT ON COLUMN material_versions.changed_by IS 'Usuario que creó esta versión';
COMMENT ON COLUMN material_versions.created_at IS 'Fecha de creación de la versión';
