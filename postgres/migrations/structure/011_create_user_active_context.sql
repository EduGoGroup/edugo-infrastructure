-- Tabla para almacenar el contexto/escuela activa del usuario
-- Permite filtrar datos en UI según la escuela seleccionada
-- Parte de FASE 1 UI Roadmap - Bloquea selector de escuela en apps

CREATE TABLE IF NOT EXISTS user_active_context (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    school_id UUID NOT NULL,
    unit_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Constraints
    CONSTRAINT uq_user_active_context_user UNIQUE(user_id),
    CONSTRAINT fk_user_active_context_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_active_context_school
        FOREIGN KEY (school_id)
        REFERENCES schools(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_active_context_unit
        FOREIGN KEY (unit_id)
        REFERENCES academic_units(id)
        ON DELETE SET NULL
);

-- Índices para performance
CREATE INDEX idx_user_active_context_user ON user_active_context(user_id);
CREATE INDEX idx_user_active_context_school ON user_active_context(school_id);

-- Trigger para updated_at automático
CREATE TRIGGER set_updated_at_user_active_context
    BEFORE UPDATE ON user_active_context
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comentarios
COMMENT ON TABLE user_active_context IS 'Almacena el contexto/escuela activa del usuario para filtrar datos en UI';
COMMENT ON COLUMN user_active_context.user_id IS 'Usuario propietario del contexto (UNIQUE: solo un contexto por usuario)';
COMMENT ON COLUMN user_active_context.school_id IS 'Escuela actualmente seleccionada por el usuario';
COMMENT ON COLUMN user_active_context.unit_id IS 'Unidad académica activa (opcional, puede ser NULL)';
COMMENT ON COLUMN user_active_context.created_at IS 'Fecha de creación del contexto';
COMMENT ON COLUMN user_active_context.updated_at IS 'Fecha de última actualización (trigger automático)';
