-- ============================================================
-- 020: iam.resources
-- Schema: iam
-- Recursos del sistema para RBAC (menú, módulos, acciones)
-- También define los ENUM types del schema iam
-- ============================================================

-- ENUM types para iam schema
CREATE TYPE iam.permission_scope AS ENUM ('system', 'school', 'unit');
CREATE TYPE iam.role_scope AS ENUM ('system', 'school', 'unit');

CREATE TABLE iam.resources (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    key character varying(50) NOT NULL,
    display_name character varying(150) NOT NULL,
    description text,
    icon character varying(100),
    parent_id uuid,
    sort_order integer DEFAULT 0 NOT NULL,
    is_menu_visible boolean DEFAULT true NOT NULL,
    scope iam.permission_scope DEFAULT 'school'::iam.permission_scope NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT resources_pkey PRIMARY KEY (id),
    CONSTRAINT resources_key_key UNIQUE (key),
    CONSTRAINT fk_resources_parent FOREIGN KEY (parent_id) REFERENCES iam.resources(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX idx_resources_key ON iam.resources USING btree (key);
CREATE INDEX idx_resources_active ON iam.resources USING btree (is_active);
CREATE INDEX idx_resources_menu_visible ON iam.resources USING btree (is_menu_visible);
CREATE INDEX idx_resources_parent ON iam.resources USING btree (parent_id);
CREATE INDEX idx_resources_sort ON iam.resources USING btree (sort_order);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON iam.resources
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
