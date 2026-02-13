-- ====================================================================
-- CONSTRAINTS: permissions
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Constraint: name debe seguir patrón resource:action
ALTER TABLE permissions ADD CONSTRAINT chk_permission_name_format
    CHECK (name ~* '^[a-z_]+:[a-z_]+(:[a-z_]+)?$');

-- Índices
CREATE INDEX idx_permissions_name ON permissions(name);
CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_scope ON permissions(scope);
CREATE INDEX idx_permissions_active ON permissions(is_active);

-- Trigger para updated_at
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
