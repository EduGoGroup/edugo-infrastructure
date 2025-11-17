-- Tabla: materials (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-mobile, worker

CREATE TABLE IF NOT EXISTS materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    uploaded_by_teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    academic_unit_id UUID REFERENCES academic_units(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    subject VARCHAR(100),
    grade VARCHAR(50),
    file_url TEXT NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'uploaded' CHECK (status IN ('uploaded', 'processing', 'ready', 'failed')),
    processing_started_at TIMESTAMP WITH TIME ZONE,
    processing_completed_at TIMESTAMP WITH TIME ZONE,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_materials_school ON materials(school_id);
CREATE INDEX idx_materials_teacher ON materials(uploaded_by_teacher_id);
CREATE INDEX idx_materials_unit ON materials(academic_unit_id);
CREATE INDEX idx_materials_status ON materials(status);
CREATE INDEX idx_materials_created_at ON materials(created_at DESC);
CREATE INDEX idx_materials_subject ON materials(subject);

COMMENT ON TABLE materials IS 'Materiales educativos subidos por docentes';
COMMENT ON COLUMN materials.status IS 'Estado: uploaded, processing, ready, failed';
