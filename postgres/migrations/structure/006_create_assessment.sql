-- Tabla: assessment (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-mobile, worker

CREATE TABLE IF NOT EXISTS assessment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL,
    mongo_document_id VARCHAR(24) NOT NULL,
    questions_count INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'generated',
    title VARCHAR(255),
    pass_threshold INTEGER DEFAULT 70,
    max_attempts INTEGER,
    time_limit_minutes INTEGER,
    total_questions INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);


COMMENT ON TABLE assessment IS 'Assessments/Quizzes generados por IA (contenido en MongoDB)';
COMMENT ON COLUMN assessment.mongo_document_id IS 'ObjectId del documento en MongoDB material_assessment';
COMMENT ON COLUMN assessment.title IS 'Título del assessment (opcional si se usa metadata de MongoDB)';
COMMENT ON COLUMN assessment.pass_threshold IS 'Porcentaje mínimo para aprobar (0-100)';
COMMENT ON COLUMN assessment.max_attempts IS 'Máximo de intentos permitidos (NULL = ilimitado)';
COMMENT ON COLUMN assessment.time_limit_minutes IS 'Límite de tiempo en minutos (NULL = sin límite)';
COMMENT ON COLUMN assessment.total_questions IS 'Total de preguntas (sincronizado con questions_count)';
