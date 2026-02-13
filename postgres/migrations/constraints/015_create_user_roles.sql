-- ====================================================================
-- CONSTRAINTS: user_roles
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Foreign Keys
ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_role
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_school
    FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_unit
    FOREIGN KEY (academic_unit_id) REFERENCES academic_units(id) ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_granted_by
    FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Unique constraint
ALTER TABLE user_roles
    ADD CONSTRAINT uq_user_role_context
    UNIQUE (user_id, role_id, school_id, academic_unit_id);

-- Check constraint: Si academic_unit_id está presente, school_id debe estar presente
ALTER TABLE user_roles
    ADD CONSTRAINT chk_user_roles_unit_requires_school
    CHECK (academic_unit_id IS NULL OR school_id IS NOT NULL);

-- Índices
CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);
CREATE INDEX idx_user_roles_school ON user_roles(school_id);
CREATE INDEX idx_user_roles_unit ON user_roles(academic_unit_id);
CREATE INDEX idx_user_roles_active ON user_roles(is_active);
CREATE INDEX idx_user_roles_expires ON user_roles(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_user_roles_user_active ON user_roles(user_id, is_active);
CREATE INDEX idx_user_roles_context ON user_roles(user_id, school_id, academic_unit_id);

-- Trigger para updated_at
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON user_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
