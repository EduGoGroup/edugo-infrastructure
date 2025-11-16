-- Tabla: assessment (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-mobile, worker

CREATE TABLE IF NOT EXISTS assessment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    mongo_document_id VARCHAR(24) NOT NULL UNIQUE,
    questions_count INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'generated' CHECK (status IN ('generated', 'published', 'archived')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_assessment_material ON assessment(material_id);
CREATE INDEX idx_assessment_mongo ON assessment(mongo_document_id);
CREATE INDEX idx_assessment_status ON assessment(status);
CREATE INDEX idx_assessment_created_at ON assessment(created_at DESC);

COMMENT ON TABLE assessment IS 'Assessments/Quizzes generados por IA (contenido en MongoDB)';
COMMENT ON COLUMN assessment.mongo_document_id IS 'ObjectId del documento en MongoDB material_assessment';
