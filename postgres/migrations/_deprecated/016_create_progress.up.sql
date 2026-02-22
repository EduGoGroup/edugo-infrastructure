-- Migration: 016_create_progress
-- Description: Crea la tabla progress para almacenar el progreso de lectura de materiales por usuario
-- Created: 2025-11-22
-- Project: edugo-infrastructure

CREATE TABLE IF NOT EXISTS progress (
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    percentage INTEGER NOT NULL DEFAULT 0 CHECK (percentage >= 0 AND percentage <= 100),
    last_page INTEGER NOT NULL DEFAULT 0 CHECK (last_page >= 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('not_started', 'in_progress', 'completed')) DEFAULT 'not_started',
    last_accessed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Primary key compuesta: un progreso por material y usuario
    PRIMARY KEY (material_id, user_id)
);

-- Índices para mejorar rendimiento
CREATE INDEX idx_progress_user_id ON progress(user_id);
CREATE INDEX idx_progress_material_id ON progress(material_id);
CREATE INDEX idx_progress_status ON progress(status);
CREATE INDEX idx_progress_last_accessed_at ON progress(last_accessed_at DESC);
CREATE INDEX idx_progress_percentage ON progress(percentage DESC);

-- Índices compuestos para consultas comunes
CREATE INDEX idx_progress_user_status ON progress(user_id, status);
CREATE INDEX idx_progress_material_status ON progress(material_id, status);

-- Comentarios para documentación
COMMENT ON TABLE progress IS 'Progreso de lectura de materiales por usuario';
COMMENT ON COLUMN progress.material_id IS 'ID del material siendo leído';
COMMENT ON COLUMN progress.user_id IS 'ID del usuario que lee el material';
COMMENT ON COLUMN progress.percentage IS 'Porcentaje de progreso (0-100)';
COMMENT ON COLUMN progress.last_page IS 'Última página leída';
COMMENT ON COLUMN progress.status IS 'Estado del progreso: not_started, in_progress, completed';
COMMENT ON COLUMN progress.last_accessed_at IS 'Última vez que se accedió al material';
COMMENT ON COLUMN progress.created_at IS 'Fecha de creación del registro de progreso';
COMMENT ON COLUMN progress.updated_at IS 'Fecha de última actualización del progreso';
