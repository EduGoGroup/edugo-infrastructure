-- ====================================================================
-- CONSTRAINTS: resources
-- VERSION: postgres/v0.17.0
-- ====================================================================

-- FK self-referencial para jerarquia
ALTER TABLE resources
    ADD CONSTRAINT fk_resources_parent
    FOREIGN KEY (parent_id) REFERENCES resources(id) ON DELETE SET NULL;

CREATE INDEX idx_resources_key ON resources(key);
CREATE INDEX idx_resources_parent ON resources(parent_id);
CREATE INDEX idx_resources_active ON resources(is_active);
CREATE INDEX idx_resources_menu_visible ON resources(is_menu_visible);
CREATE INDEX idx_resources_sort ON resources(sort_order);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
