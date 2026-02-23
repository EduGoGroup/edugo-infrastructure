-- ============================================================
-- 060: ui_config.screen_templates
-- Schema: ui_config
-- Plantillas de pantalla (patrones reutilizables de UI)
-- ============================================================

CREATE TABLE ui_config.screen_templates (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    pattern character varying(50) NOT NULL,
    name character varying(200) NOT NULL,
    description text,
    version integer DEFAULT 1 NOT NULL,
    definition jsonb NOT NULL,
    is_active boolean DEFAULT true,
    created_by uuid,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    CONSTRAINT screen_templates_pkey PRIMARY KEY (id),
    CONSTRAINT screen_templates_name_version_key UNIQUE (name, version)
);

-- Indexes
CREATE INDEX idx_screen_templates_pattern ON ui_config.screen_templates USING btree (pattern);
CREATE INDEX idx_screen_templates_active ON ui_config.screen_templates USING btree (is_active) WHERE (is_active = true);
CREATE INDEX idx_screen_templates_definition ON ui_config.screen_templates USING gin (definition);

CREATE TRIGGER update_screen_templates_updated_at BEFORE UPDATE ON ui_config.screen_templates
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
