-- FK
ALTER TABLE ui_config.screen_templates
    ADD CONSTRAINT fk_screen_templates_created_by
    FOREIGN KEY (created_by) REFERENCES public.users(id);

-- Índices
CREATE INDEX idx_screen_templates_pattern ON ui_config.screen_templates(pattern);
CREATE INDEX idx_screen_templates_active ON ui_config.screen_templates(is_active) WHERE is_active = true;
CREATE INDEX idx_screen_templates_definition ON ui_config.screen_templates USING GIN (definition);

-- Trigger updated_at
CREATE TRIGGER update_screen_templates_updated_at
    BEFORE UPDATE ON ui_config.screen_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
