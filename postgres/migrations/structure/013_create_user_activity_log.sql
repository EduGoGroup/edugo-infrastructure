-- Tipo ENUM para clasificar actividades del usuario
CREATE TYPE activity_type AS ENUM (
    'material_started',      -- Usuario inició un material
    'material_progress',     -- Usuario avanzó en lectura
    'material_completed',    -- Usuario completó material
    'summary_viewed',        -- Usuario vio resumen generado
    'quiz_started',          -- Usuario inició quiz
    'quiz_completed',        -- Usuario completó quiz
    'quiz_passed',          -- Usuario aprobó quiz
    'quiz_failed'           -- Usuario reprobó quiz
);

-- Tabla para log de actividades del usuario
-- Uso: historial, analytics, actividad reciente en Home
-- Parte de FASE 1 UI Roadmap - Actividad reciente en apps

CREATE TABLE IF NOT EXISTS user_activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    activity_type activity_type NOT NULL,
    material_id UUID,
    school_id UUID,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user_activity_log_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_activity_log_material
        FOREIGN KEY (material_id)
        REFERENCES materials(id)
        ON DELETE SET NULL,
    CONSTRAINT fk_user_activity_log_school
        FOREIGN KEY (school_id)
        REFERENCES schools(id)
        ON DELETE SET NULL
);

-- Índices para queries frecuentes
CREATE INDEX idx_user_activity_user_created
    ON user_activity_log(user_id, created_at DESC);

CREATE INDEX idx_user_activity_school
    ON user_activity_log(school_id, created_at DESC);

CREATE INDEX idx_user_activity_type
    ON user_activity_log(activity_type);

-- Índice para rate limiting
-- Nota: No se puede usar índice parcial con NOW() porque no es IMMUTABLE
-- El índice completo es suficiente para queries de rate limiting
CREATE INDEX idx_user_activity_rate_limit
    ON user_activity_log(user_id, activity_type, created_at);

-- Comentarios
COMMENT ON TABLE user_activity_log IS 'Log de actividades del usuario para historial y analytics';
COMMENT ON COLUMN user_activity_log.activity_type IS 'Tipo de actividad realizada';
COMMENT ON COLUMN user_activity_log.material_id IS 'Material asociado (NULL si no aplica, SET NULL si material eliminado)';
COMMENT ON COLUMN user_activity_log.school_id IS 'Escuela asociada (NULL si no aplica, SET NULL si escuela eliminada)';
COMMENT ON COLUMN user_activity_log.metadata IS 'Datos adicionales en JSON (ej: score, pages, time_spent)';
COMMENT ON COLUMN user_activity_log.created_at IS 'Timestamp de la actividad';
