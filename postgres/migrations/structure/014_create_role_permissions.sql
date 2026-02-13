-- ====================================================================
-- TABLA: role_permissions
-- DESCRIPCIÓN: Relación N:N entre roles y permisos
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON TABLE role_permissions IS 'Relación N:N entre roles y permisos (RBAC)';
