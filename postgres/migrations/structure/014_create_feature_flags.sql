-- Tabla para sistema de Feature Flags
-- Control remoto de características en apps sin redeployar
-- Deuda Técnica - Requerido por Apple App
-- Spec: /Users/jhoanmedina/source/EduGo/EduUI/apple-app/docs/backend-specs/feature-flags/BACKEND-SPEC-FEATURE-FLAGS.md

CREATE TABLE IF NOT EXISTS feature_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Identificación
    key VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,

    -- Estado
    enabled BOOLEAN NOT NULL DEFAULT false,
    enabled_globally BOOLEAN NOT NULL DEFAULT false,

    -- Metadata
    category VARCHAR(50),
    priority INTEGER NOT NULL DEFAULT 0,

    -- Restricciones de versión
    minimum_app_version VARCHAR(20),
    minimum_build_number INTEGER,
    maximum_app_version VARCHAR(20),
    maximum_build_number INTEGER,

    -- Segmentación (Phase 2)
    enabled_for_roles JSONB DEFAULT '[]',
    enabled_for_user_ids JSONB DEFAULT '[]',
    disabled_for_user_ids JSONB DEFAULT '[]',

    -- A/B Testing (Phase 3)
    rollout_percentage INTEGER DEFAULT 100,

    -- Flags de control
    is_experimental BOOLEAN NOT NULL DEFAULT false,
    requires_restart BOOLEAN NOT NULL DEFAULT false,
    is_debug_only BOOLEAN NOT NULL DEFAULT false,
    affects_security BOOLEAN NOT NULL DEFAULT false,

    -- Auditoría
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID,
    updated_by UUID,

    -- Constraints
    CONSTRAINT fk_feature_flags_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE SET NULL,
    CONSTRAINT fk_feature_flags_updated_by
        FOREIGN KEY (updated_by)
        REFERENCES users(id)
        ON DELETE SET NULL,
    CONSTRAINT chk_valid_rollout
        CHECK (rollout_percentage >= 0 AND rollout_percentage <= 100),
    CONSTRAINT chk_valid_build_numbers
        CHECK (
            minimum_build_number IS NULL
            OR maximum_build_number IS NULL
            OR minimum_build_number <= maximum_build_number
        )
);

-- Índices para performance
CREATE INDEX idx_feature_flags_key ON feature_flags(key);
CREATE INDEX idx_feature_flags_enabled ON feature_flags(enabled);
CREATE INDEX idx_feature_flags_category ON feature_flags(category);
CREATE INDEX idx_feature_flags_updated_at ON feature_flags(updated_at);

-- Trigger para updated_at automático
CREATE TRIGGER set_updated_at_feature_flags
    BEFORE UPDATE ON feature_flags
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentarios
COMMENT ON TABLE feature_flags IS 'Sistema de feature flags para control remoto de características en apps';
COMMENT ON COLUMN feature_flags.key IS 'Identificador único del feature flag (ej: biometric_login)';
COMMENT ON COLUMN feature_flags.enabled IS 'Estado global del feature flag';
COMMENT ON COLUMN feature_flags.enabled_globally IS 'Si true, ignora segmentación y habilita para todos';
COMMENT ON COLUMN feature_flags.rollout_percentage IS 'Porcentaje de usuarios habilitados (0-100) para A/B testing';
COMMENT ON COLUMN feature_flags.enabled_for_roles IS 'Array JSON de roles habilitados (ej: ["admin", "teacher"])';
COMMENT ON COLUMN feature_flags.enabled_for_user_ids IS 'Array JSON de UUIDs de usuarios habilitados (whitelist)';
COMMENT ON COLUMN feature_flags.disabled_for_user_ids IS 'Array JSON de UUIDs de usuarios deshabilitados (blacklist)';
COMMENT ON COLUMN feature_flags.is_experimental IS 'Marca feature como experimental en UI';
COMMENT ON COLUMN feature_flags.requires_restart IS 'Si true, requiere reiniciar app para aplicar cambio';
COMMENT ON COLUMN feature_flags.is_debug_only IS 'Si true, solo habilitado en builds de debug/desarrollo';
COMMENT ON COLUMN feature_flags.affects_security IS 'Si true, cambios requieren aprobación especial';
