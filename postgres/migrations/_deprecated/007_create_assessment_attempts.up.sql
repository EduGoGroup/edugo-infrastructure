-- Tabla: assessment_attempt (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-mobile

CREATE TABLE IF NOT EXISTS assessment_attempt (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessment(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    score DECIMAL(5,2),
    max_score DECIMAL(5,2),
    percentage DECIMAL(5,2),
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'abandoned')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attempt_assessment ON assessment_attempt(assessment_id);
CREATE INDEX idx_attempt_student ON assessment_attempt(student_id);
CREATE INDEX idx_attempt_status ON assessment_attempt(status);
CREATE INDEX idx_attempt_completed_at ON assessment_attempt(completed_at DESC);

COMMENT ON TABLE assessment_attempt IS 'Intentos de estudiantes en assessments';
