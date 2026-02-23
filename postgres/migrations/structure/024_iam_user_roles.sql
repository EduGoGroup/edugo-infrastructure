-- ============================================================
-- 024: iam.user_roles
-- Schema: iam
-- Asignación de roles a usuarios con contexto (school/unit)
-- Cross-schema FKs (user_id, school_id, academic_unit_id, granted_by) van en 070
-- ============================================================

CREATE TABLE iam.user_roles (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    school_id uuid,
    academic_unit_id uuid,
    is_active boolean DEFAULT true NOT NULL,
    granted_by uuid,
    granted_at timestamptz DEFAULT now() NOT NULL,
    expires_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT user_roles_pkey PRIMARY KEY (id),
    CONSTRAINT uq_user_role_context UNIQUE (user_id, role_id, school_id, academic_unit_id),
    CONSTRAINT chk_user_roles_unit_requires_school CHECK (academic_unit_id IS NULL OR school_id IS NOT NULL),
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES iam.roles(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_user_roles_user ON iam.user_roles USING btree (user_id);
CREATE INDEX idx_user_roles_role ON iam.user_roles USING btree (role_id);
CREATE INDEX idx_user_roles_school ON iam.user_roles USING btree (school_id);
CREATE INDEX idx_user_roles_unit ON iam.user_roles USING btree (academic_unit_id);
CREATE INDEX idx_user_roles_active ON iam.user_roles USING btree (is_active);
CREATE INDEX idx_user_roles_context ON iam.user_roles USING btree (user_id, school_id, academic_unit_id);
CREATE INDEX idx_user_roles_user_active ON iam.user_roles USING btree (user_id, is_active);
CREATE INDEX idx_user_roles_expires ON iam.user_roles USING btree (expires_at) WHERE (expires_at IS NOT NULL);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON iam.user_roles
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
