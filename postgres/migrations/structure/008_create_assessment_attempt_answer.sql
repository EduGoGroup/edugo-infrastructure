-- Tabla: assessment_attempt_answer (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-mobile

CREATE TABLE IF NOT EXISTS assessment_attempt_answer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL,
    question_index INTEGER NOT NULL,
    student_answer TEXT,
    is_correct BOOLEAN,
    points_earned DECIMAL(5,2),
    max_points DECIMAL(5,2),
    time_spent_seconds INTEGER,
    answered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


COMMENT ON TABLE assessment_attempt_answer IS 'Respuestas individuales de estudiantes por pregunta';
COMMENT ON COLUMN assessment_attempt_answer.question_index IS 'Índice de la pregunta (0-based). APIs mapean a question_id según necesidad.';
COMMENT ON COLUMN assessment_attempt_answer.student_answer IS 'Respuesta del estudiante (TEXT flexible: JSON, string, etc). APIs mapean a selected_answer_id según necesidad.';
COMMENT ON COLUMN assessment_attempt_answer.time_spent_seconds IS 'Tiempo que tomó responder esta pregunta en segundos';
