-- ============================================================
-- 090: academic.grades
-- Schema: academic
-- Calificaciones de estudiantes por materia y periodo
-- ============================================================

CREATE TABLE academic.grades (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    membership_id uuid NOT NULL,
    subject_id uuid NOT NULL,
    period_id uuid NOT NULL,
    grade_value decimal(5,2),
    grade_letter character varying(5),
    assessment_scores jsonb DEFAULT '[]'::jsonb,
    teacher_id uuid,
    notes text,
    finalized_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT grades_pkey PRIMARY KEY (id),
    CONSTRAINT grades_unique UNIQUE (membership_id, subject_id, period_id),
    CONSTRAINT grades_membership_fkey FOREIGN KEY (membership_id) REFERENCES academic.memberships(id) ON DELETE CASCADE,
    CONSTRAINT grades_subject_fkey FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE CASCADE,
    CONSTRAINT grades_period_fkey FOREIGN KEY (period_id) REFERENCES academic.academic_periods(id) ON DELETE CASCADE,
    CONSTRAINT grades_teacher_fkey FOREIGN KEY (teacher_id) REFERENCES academic.memberships(id)
);

CREATE INDEX idx_grades_membership ON academic.grades(membership_id);
CREATE INDEX idx_grades_period ON academic.grades(period_id);
CREATE INDEX idx_grades_subject ON academic.grades(subject_id);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.grades
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
