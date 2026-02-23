-- ============================================================
-- 021: iam.roles
-- Schema: iam
-- Roles del sistema RBAC
-- ============================================================

CREATE TABLE iam.roles (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name character varying(50) NOT NULL,
    display_name character varying(100) NOT NULL,
    description text,
    scope iam.role_scope DEFAULT 'school'::iam.role_scope NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT roles_pkey PRIMARY KEY (id),
    CONSTRAINT roles_name_key UNIQUE (name)
);

-- Indexes
CREATE INDEX idx_roles_name ON iam.roles USING btree (name);
CREATE INDEX idx_roles_active ON iam.roles USING btree (is_active);
CREATE INDEX idx_roles_scope ON iam.roles USING btree (scope);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON iam.roles
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
