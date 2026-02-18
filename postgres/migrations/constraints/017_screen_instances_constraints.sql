-- FK
ALTER TABLE ui_config.screen_instances
    ADD CONSTRAINT fk_screen_instances_template
    FOREIGN KEY (template_id) REFERENCES ui_config.screen_templates(id);

ALTER TABLE ui_config.screen_instances
    ADD CONSTRAINT fk_screen_instances_created_by
    FOREIGN KEY (created_by) REFERENCES public.users(id);

-- Índices
CREATE INDEX idx_screen_instances_template ON ui_config.screen_instances(template_id);
CREATE INDEX idx_screen_instances_scope ON ui_config.screen_instances(scope);
CREATE INDEX idx_screen_instances_active ON ui_config.screen_instances(is_active) WHERE is_active = true;
CREATE INDEX idx_screen_instances_slot_data ON ui_config.screen_instances USING GIN (slot_data);

CREATE INDEX idx_screen_instances_handler_key ON ui_config.screen_instances(handler_key) WHERE handler_key IS NOT NULL;

-- Trigger updated_at
CREATE TRIGGER update_screen_instances_updated_at
    BEFORE UPDATE ON ui_config.screen_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
