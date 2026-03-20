-- ============================================================
-- 061: ui_config.screen_instances
-- Schema: ui_config
-- Instancias concretas de pantalla basadas en plantillas
-- ============================================================

CREATE TABLE ui_config.screen_instances (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    screen_key character varying(100) NOT NULL,
    template_id uuid NOT NULL,
    name character varying(200) NOT NULL,
    description text,
    slot_data jsonb DEFAULT '{}'::jsonb NOT NULL,
    scope character varying(20) DEFAULT 'school',
    required_permission character varying(100),
    handler_key character varying(100),
    is_active boolean DEFAULT true,
    created_by uuid,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    CONSTRAINT screen_instances_pkey PRIMARY KEY (id),
    CONSTRAINT screen_instances_screen_key_key UNIQUE (screen_key),
    -- Intra-schema FK
    CONSTRAINT fk_screen_instances_template FOREIGN KEY (template_id) REFERENCES ui_config.screen_templates(id)
);

-- Indexes
CREATE INDEX idx_screen_instances_template ON ui_config.screen_instances USING btree (template_id);
CREATE INDEX idx_screen_instances_active ON ui_config.screen_instances USING btree (is_active) WHERE (is_active = true);
CREATE INDEX idx_screen_instances_scope ON ui_config.screen_instances USING btree (scope);
CREATE INDEX idx_screen_instances_slot_data ON ui_config.screen_instances USING gin (slot_data);
CREATE INDEX idx_screen_instances_handler_key ON ui_config.screen_instances USING btree (handler_key) WHERE (handler_key IS NOT NULL);

CREATE TRIGGER update_screen_instances_updated_at BEFORE UPDATE ON ui_config.screen_instances
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
