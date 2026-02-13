-- ====================================================================
-- CONSTRAINTS: role_permissions
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Foreign Keys
ALTER TABLE role_permissions
    ADD CONSTRAINT fk_role_permissions_role
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

ALTER TABLE role_permissions
    ADD CONSTRAINT fk_role_permissions_permission
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;

-- Unique constraint: Un rol no puede tener el mismo permiso duplicado
ALTER TABLE role_permissions
    ADD CONSTRAINT uq_role_permission UNIQUE (role_id, permission_id);

-- Índices
CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);
