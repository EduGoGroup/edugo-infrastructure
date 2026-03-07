-- ============================================================
-- 053: assessment.assessment_materials
-- Schema: assessment
-- Relacion N:N entre assessments y materiales
-- Cross-schema FK (material_id -> content.materials) inline
-- ============================================================

CREATE TABLE IF NOT EXISTS assessment.assessment_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL,
    material_id UUID NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT assessment_materials_unique UNIQUE(assessment_id, material_id),
    CONSTRAINT assessment_materials_assessment_fk FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE,
    CONSTRAINT assessment_materials_material_fk FOREIGN KEY (material_id) REFERENCES content.materials(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_assessment_materials_assessment ON assessment.assessment_materials(assessment_id);
CREATE INDEX IF NOT EXISTS idx_assessment_materials_material ON assessment.assessment_materials(material_id);
