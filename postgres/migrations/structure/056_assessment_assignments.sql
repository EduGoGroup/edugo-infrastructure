-- ============================================================
-- 056: assessment.assessment_assignments
-- Schema: assessment
-- Asignaciones de evaluaciones a estudiantes o unidades academicas
-- Cross-schema FKs (student_id, assigned_by -> auth.users,
--   academic_unit_id -> academic.academic_units) van en 070
-- ============================================================

CREATE TABLE assessment.assessment_assignments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    assessment_id uuid NOT NULL,
    student_id uuid,
    academic_unit_id uuid,
    assigned_by uuid NOT NULL,
    assigned_at timestamptz DEFAULT now() NOT NULL,
    due_date timestamptz,
    CONSTRAINT assessment_assignments_pkey PRIMARY KEY (id),
    CONSTRAINT chk_assignment_target CHECK (
        (student_id IS NOT NULL AND academic_unit_id IS NULL) OR
        (student_id IS NULL AND academic_unit_id IS NOT NULL)
    ),
    -- Intra-schema FK
    CONSTRAINT assessment_assignments_assessment_fkey FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE
);

CREATE INDEX idx_assignment_assessment ON assessment.assessment_assignments(assessment_id);
CREATE INDEX idx_assignment_student ON assessment.assessment_assignments(student_id) WHERE student_id IS NOT NULL;
CREATE INDEX idx_assignment_unit ON assessment.assessment_assignments(academic_unit_id) WHERE academic_unit_id IS NOT NULL;
CREATE UNIQUE INDEX idx_unique_student_assignment ON assessment.assessment_assignments(assessment_id, student_id) WHERE student_id IS NOT NULL;
CREATE UNIQUE INDEX idx_unique_unit_assignment ON assessment.assessment_assignments(assessment_id, academic_unit_id) WHERE academic_unit_id IS NOT NULL;
