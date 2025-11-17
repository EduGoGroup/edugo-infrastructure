-- Tabla: assessment_attempt_answer (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-mobile

CREATE TABLE IF NOT EXISTS assessment_attempt_answer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL REFERENCES assessment_attempt(id) ON DELETE CASCADE,
    question_index INTEGER NOT NULL,
    student_answer TEXT,
    is_correct BOOLEAN,
    points_earned DECIMAL(5,2),
    max_points DECIMAL(5,2),
    answered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(attempt_id, question_index)
);

CREATE INDEX idx_answer_attempt ON assessment_attempt_answer(attempt_id);
CREATE INDEX idx_answer_correct ON assessment_attempt_answer(is_correct);

COMMENT ON TABLE assessment_attempt_answer IS 'Respuestas individuales de estudiantes por pregunta';
COMMENT ON COLUMN assessment_attempt_answer.question_index IS 'Índice de la pregunta en el assessment (0-based)';
