-- ====================================================================
-- TABLA: user_roles
-- DESCRIPCIÓN: Asignación de roles a usuarios en contextos específicos
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    school_id UUID,
    academic_unit_id UUID,
    is_active BOOLEAN DEFAULT true NOT NULL,
    granted_by UUID,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON TABLE user_roles IS 'Asignación de roles a usuarios en contextos específicos (RBAC)';
COMMENT ON COLUMN user_roles.school_id IS 'Escuela en la que aplica el rol. NULL = rol a nivel sistema';
COMMENT ON COLUMN user_roles.academic_unit_id IS 'Unidad académica en la que aplica el rol. NULL = rol a nivel escuela';
COMMENT ON COLUMN user_roles.granted_by IS 'Usuario que otorgó el rol (auditoría)';
COMMENT ON COLUMN user_roles.expires_at IS 'Fecha de expiración del rol. NULL = no expira';
