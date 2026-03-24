-- ============================================================
-- 092: academic.schedules
-- Schema: academic
-- Horarios semanales por unidad academica
-- ============================================================

CREATE TABLE academic.schedules (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    academic_unit_id uuid NOT NULL,
    subject_id uuid NOT NULL,
    teacher_membership_id uuid NOT NULL,
    day_of_week integer NOT NULL,
    start_time time NOT NULL,
    end_time time NOT NULL,
    room character varying(50),
    period_id uuid,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT schedules_pkey PRIMARY KEY (id),
    CONSTRAINT schedules_dow_check CHECK (day_of_week BETWEEN 0 AND 6),
    CONSTRAINT schedules_unit_fkey FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE CASCADE,
    CONSTRAINT schedules_subject_fkey FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE CASCADE,
    CONSTRAINT schedules_teacher_fkey FOREIGN KEY (teacher_membership_id) REFERENCES academic.memberships(id) ON DELETE CASCADE,
    CONSTRAINT schedules_period_fkey FOREIGN KEY (period_id) REFERENCES academic.academic_periods(id)
);

CREATE INDEX idx_schedules_unit ON academic.schedules(academic_unit_id);
CREATE INDEX idx_schedules_teacher ON academic.schedules(teacher_membership_id);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.schedules
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
