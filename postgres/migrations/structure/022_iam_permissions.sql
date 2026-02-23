-- ============================================================
-- 022: iam.permissions
-- Schema: iam
-- Permisos granulares del sistema RBAC
-- ============================================================

CREATE TABLE iam.permissions (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    display_name character varying(150) NOT NULL,
    description text,
    resource_id uuid NOT NULL,
    action character varying(50) NOT NULL,
    scope iam.permission_scope DEFAULT 'school'::iam.permission_scope NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT permissions_pkey PRIMARY KEY (id),
    CONSTRAINT permissions_name_key UNIQUE (name),
    CONSTRAINT uq_permissions_resource_action UNIQUE (resource_id, action),
    CONSTRAINT chk_permission_name_format CHECK (name ~* '^[a-z_]+:[a-z_]+(:[a-z_]+)?$'),
    CONSTRAINT fk_permissions_resource FOREIGN KEY (resource_id) REFERENCES iam.resources(id) ON DELETE RESTRICT
);

-- Indexes
CREATE INDEX idx_permissions_name ON iam.permissions USING btree (name);
CREATE INDEX idx_permissions_resource ON iam.permissions USING btree (resource_id);
CREATE INDEX idx_permissions_active ON iam.permissions USING btree (is_active);
CREATE INDEX idx_permissions_scope ON iam.permissions USING btree (scope);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON iam.permissions
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
