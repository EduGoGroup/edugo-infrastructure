-- Migration: 015_create_guardian_relations
-- Description: Crea la tabla guardian_relations para almacenar relaciones entre apoderados y estudiantes
-- Created: 2025-11-22
-- Project: edugo-infrastructure

CREATE TABLE IF NOT EXISTS guardian_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    guardian_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    relationship_type VARCHAR(50) NOT NULL CHECK (relationship_type IN ('father', 'mother', 'grandfather', 'grandmother', 'uncle', 'aunt', 'sibling', 'legal_guardian', 'other')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,

    -- Constraints
    UNIQUE(guardian_id, student_id),
    CHECK (guardian_id != student_id) -- El apoderado no puede ser el mismo estudiante
);

-- Índices para mejorar rendimiento
CREATE INDEX idx_guardian_relations_guardian_id ON guardian_relations(guardian_id);
CREATE INDEX idx_guardian_relations_student_id ON guardian_relations(student_id);
CREATE INDEX idx_guardian_relations_relationship_type ON guardian_relations(relationship_type);
CREATE INDEX idx_guardian_relations_is_active ON guardian_relations(is_active);
CREATE INDEX idx_guardian_relations_created_at ON guardian_relations(created_at DESC);

-- Índice compuesto para consultas comunes
CREATE INDEX idx_guardian_relations_active_guardian ON guardian_relations(guardian_id, is_active);
CREATE INDEX idx_guardian_relations_active_student ON guardian_relations(student_id, is_active);

-- Comentarios para documentación
COMMENT ON TABLE guardian_relations IS 'Relaciones entre apoderados (guardians) y estudiantes';
COMMENT ON COLUMN guardian_relations.id IS 'Identificador único de la relación';
COMMENT ON COLUMN guardian_relations.guardian_id IS 'ID del usuario que actúa como apoderado';
COMMENT ON COLUMN guardian_relations.student_id IS 'ID del usuario que es el estudiante';
COMMENT ON COLUMN guardian_relations.relationship_type IS 'Tipo de relación familiar o legal';
COMMENT ON COLUMN guardian_relations.is_active IS 'Indica si la relación está activa';
COMMENT ON COLUMN guardian_relations.created_at IS 'Fecha de creación de la relación';
COMMENT ON COLUMN guardian_relations.updated_at IS 'Fecha de última actualización';
COMMENT ON COLUMN guardian_relations.created_by IS 'Usuario que creó la relación';
