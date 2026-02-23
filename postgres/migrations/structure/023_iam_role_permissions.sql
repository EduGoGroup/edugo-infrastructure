-- ============================================================
-- 023: iam.role_permissions
-- Schema: iam
-- Tabla de unión entre roles y permisos
-- ============================================================

CREATE TABLE iam.role_permissions (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    role_id uuid NOT NULL,
    permission_id uuid NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT role_permissions_pkey PRIMARY KEY (id),
    CONSTRAINT uq_role_permission UNIQUE (role_id, permission_id),
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES iam.roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES iam.permissions(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_role_permissions_role ON iam.role_permissions USING btree (role_id);
CREATE INDEX idx_role_permissions_permission ON iam.role_permissions USING btree (permission_id);
