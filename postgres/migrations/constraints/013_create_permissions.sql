-- ====================================================================
-- CONSTRAINTS: permissions
-- VERSION: postgres/v0.17.0
-- ====================================================================

ALTER TABLE permissions
    ADD CONSTRAINT fk_permissions_resource
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE RESTRICT;

ALTER TABLE permissions
    ADD CONSTRAINT chk_permission_name_format
    CHECK (name ~* '^[a-z_]+:[a-z_]+(:[a-z_]+)?$');

ALTER TABLE permissions
    ADD CONSTRAINT uq_permissions_resource_action
    UNIQUE (resource_id, action);

CREATE INDEX idx_permissions_name ON permissions(name);
CREATE INDEX idx_permissions_resource ON permissions(resource_id);
CREATE INDEX idx_permissions_scope ON permissions(scope);
CREATE INDEX idx_permissions_active ON permissions(is_active);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
