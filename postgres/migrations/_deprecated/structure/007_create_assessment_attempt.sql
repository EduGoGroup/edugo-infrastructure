-- Tabla: assessment_attempt (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-mobile

CREATE TABLE IF NOT EXISTS assessment_attempt (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL,
    student_id UUID NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    score DECIMAL(5,2),
    max_score DECIMAL(5,2),
    percentage DECIMAL(5,2),
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress',
    time_spent_seconds INTEGER,
    idempotency_key VARCHAR(64),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE assessment_attempt IS 'Intentos de estudiantes en assessments';
COMMENT ON COLUMN assessment_attempt.time_spent_seconds IS 'Tiempo total del intento en segundos (max 2 horas)';
COMMENT ON COLUMN assessment_attempt.idempotency_key IS 'Clave para prevenir intentos duplicados';
