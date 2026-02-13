-- ====================================================================
-- CONSTRAINTS: roles
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Índices
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_scope ON roles(scope);
CREATE INDEX idx_roles_active ON roles(is_active);

-- Trigger para updated_at (usa función existente update_updated_at_column)
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
