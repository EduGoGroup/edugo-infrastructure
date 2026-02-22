-- Tabla: academic_units (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-admin (jerarquía), api-mobile (plano), worker
-- Versión: v0.7.0 (extendida con jerarquía para api-admin)

CREATE TABLE IF NOT EXISTS academic_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_unit_id UUID,
    school_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    level VARCHAR(50),
    academic_year INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);


COMMENT ON TABLE academic_units IS 'Unidades académicas con soporte de jerarquía opcional';
COMMENT ON COLUMN academic_units.parent_unit_id IS 'Unidad padre (jerarquía: Facultad → Departamento). NULL = raíz';
COMMENT ON COLUMN academic_units.type IS 'Tipo: school, grade, class, section, club, department';
COMMENT ON COLUMN academic_units.description IS 'Descripción de la unidad académica';
COMMENT ON COLUMN academic_units.academic_year IS 'Año académico. 0 = sin año específico (usado por api-admin)';
COMMENT ON COLUMN academic_units.metadata IS 'Metadata extensible';
