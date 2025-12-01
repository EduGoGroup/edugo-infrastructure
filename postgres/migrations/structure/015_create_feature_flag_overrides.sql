-- Tabla para sobrescrituras específicas de feature flags por usuario
-- Permite habilitar/deshabilitar features para usuarios individuales
-- Deuda Técnica - Phase 2 de Feature Flags
-- Spec: /Users/jhoanmedina/source/EduGo/EduUI/apple-app/docs/backend-specs/feature-flags/BACKEND-SPEC-FEATURE-FLAGS.md

CREATE TABLE IF NOT EXISTS feature_flag_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feature_flag_id UUID NOT NULL,
    user_id UUID NOT NULL,
    enabled BOOLEAN NOT NULL,
    reason TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID,

    -- Constraints
    CONSTRAINT fk_ff_overrides_flag
        FOREIGN KEY (feature_flag_id)
        REFERENCES feature_flags(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_ff_overrides_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_ff_overrides_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE SET NULL,
    CONSTRAINT uq_ff_overrides_flag_user
        UNIQUE(feature_flag_id, user_id)
);

-- Índices para performance
CREATE INDEX idx_ff_overrides_user ON feature_flag_overrides(user_id);
CREATE INDEX idx_ff_overrides_flag ON feature_flag_overrides(feature_flag_id);
CREATE INDEX idx_ff_overrides_expires ON feature_flag_overrides(expires_at)
    WHERE expires_at IS NOT NULL;

-- Comentarios
COMMENT ON TABLE feature_flag_overrides IS 'Sobrescrituras de feature flags específicas por usuario (opcional, temporal)';
COMMENT ON COLUMN feature_flag_overrides.enabled IS 'Estado override para este usuario (sobrescribe lógica global)';
COMMENT ON COLUMN feature_flag_overrides.reason IS 'Razón de la sobrescritura (debugging, testing, etc.)';
COMMENT ON COLUMN feature_flag_overrides.expires_at IS 'Fecha de expiración del override (NULL = permanente)';
COMMENT ON COLUMN feature_flag_overrides.created_by IS 'Admin que creó el override';
