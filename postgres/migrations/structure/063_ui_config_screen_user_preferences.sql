-- ============================================================
-- 063: ui_config.screen_user_preferences
-- Schema: ui_config
-- Preferencias de usuario por instancia de pantalla
-- Cross-schema FK (user_id -> auth.users) va en 070
-- ============================================================

CREATE TABLE ui_config.screen_user_preferences (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    screen_instance_id uuid NOT NULL,
    user_id uuid NOT NULL,
    preferences jsonb DEFAULT '{}'::jsonb NOT NULL,
    updated_at timestamptz DEFAULT now(),
    CONSTRAINT screen_user_preferences_pkey PRIMARY KEY (id),
    CONSTRAINT screen_user_preferences_screen_instance_id_user_id_key UNIQUE (screen_instance_id, user_id),
    -- Intra-schema FK
    CONSTRAINT fk_screen_user_prefs_instance FOREIGN KEY (screen_instance_id) REFERENCES ui_config.screen_instances(id)
);

-- Indexes
CREATE INDEX idx_screen_user_prefs_screen ON ui_config.screen_user_preferences USING btree (screen_instance_id);
CREATE INDEX idx_screen_user_prefs_user ON ui_config.screen_user_preferences USING btree (user_id);

CREATE TRIGGER update_screen_user_prefs_updated_at BEFORE UPDATE ON ui_config.screen_user_preferences
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
