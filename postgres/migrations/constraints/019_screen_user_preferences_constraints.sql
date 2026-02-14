-- FK
ALTER TABLE ui_config.screen_user_preferences
    ADD CONSTRAINT fk_screen_user_prefs_instance
    FOREIGN KEY (screen_instance_id) REFERENCES ui_config.screen_instances(id);

ALTER TABLE ui_config.screen_user_preferences
    ADD CONSTRAINT fk_screen_user_prefs_user
    FOREIGN KEY (user_id) REFERENCES public.users(id);

-- Índices
CREATE INDEX idx_screen_user_prefs_user ON ui_config.screen_user_preferences(user_id);
CREATE INDEX idx_screen_user_prefs_screen ON ui_config.screen_user_preferences(screen_instance_id);

-- Trigger updated_at
CREATE TRIGGER update_screen_user_prefs_updated_at
    BEFORE UPDATE ON ui_config.screen_user_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
