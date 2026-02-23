-- ============================================================
-- 062: ui_config.resource_screens
-- Schema: ui_config
-- Asociación entre recursos IAM y pantallas
-- Cross-schema FK (resource_id -> iam.resources) va en 070
-- ============================================================

CREATE TABLE ui_config.resource_screens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    resource_id uuid NOT NULL,
    resource_key character varying(100) NOT NULL,
    screen_key character varying(100) NOT NULL,
    screen_type character varying(50) NOT NULL,
    is_default boolean DEFAULT false,
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    CONSTRAINT resource_screens_pkey PRIMARY KEY (id),
    CONSTRAINT resource_screens_resource_id_screen_type_key UNIQUE (resource_id, screen_type),
    -- Intra-schema FK
    CONSTRAINT fk_resource_screens_screen_key FOREIGN KEY (screen_key) REFERENCES ui_config.screen_instances(screen_key)
);

-- Indexes
CREATE INDEX idx_resource_screens_resource ON ui_config.resource_screens USING btree (resource_id);
CREATE INDEX idx_resource_screens_resource_key ON ui_config.resource_screens USING btree (resource_key);
CREATE INDEX idx_resource_screens_screen_key ON ui_config.resource_screens USING btree (screen_key);

CREATE TRIGGER update_resource_screens_updated_at BEFORE UPDATE ON ui_config.resource_screens
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
