-- FK
ALTER TABLE ui_config.resource_screens
    ADD CONSTRAINT fk_resource_screens_resource
    FOREIGN KEY (resource_id) REFERENCES public.resources(id);

ALTER TABLE ui_config.resource_screens
    ADD CONSTRAINT fk_resource_screens_screen_key
    FOREIGN KEY (screen_key) REFERENCES ui_config.screen_instances(screen_key);

-- Índices
CREATE INDEX idx_resource_screens_resource ON ui_config.resource_screens(resource_id);
CREATE INDEX idx_resource_screens_resource_key ON ui_config.resource_screens(resource_key);
CREATE INDEX idx_resource_screens_screen_key ON ui_config.resource_screens(screen_key);

-- Trigger updated_at
CREATE TRIGGER update_resource_screens_updated_at
    BEFORE UPDATE ON ui_config.resource_screens
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
