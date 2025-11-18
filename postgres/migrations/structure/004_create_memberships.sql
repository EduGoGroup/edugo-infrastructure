-- Tabla: memberships (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-admin, api-mobile, worker
-- Versión: v0.7.0 (extendida con roles administrativos para api-admin)

CREATE TABLE IF NOT EXISTS memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    school_id UUID NOT NULL,
    academic_unit_id UUID,
    role VARCHAR(50) NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT true,
    enrolled_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    withdrawn_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


COMMENT ON TABLE memberships IS 'Relación usuario-escuela-unidad académica';
COMMENT ON COLUMN memberships.role IS 'Rol: teacher, student, guardian, coordinator, admin, assistant';
COMMENT ON COLUMN memberships.metadata IS 'Metadata extensible: permisos específicos, configuración, historial';
COMMENT ON COLUMN memberships.enrolled_at IS 'Fecha de inicio de membresía';
COMMENT ON COLUMN memberships.withdrawn_at IS 'Fecha de fin de membresía (NULL = activo)';
