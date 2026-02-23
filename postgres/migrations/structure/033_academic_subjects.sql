-- ============================================================
-- 033: academic.subjects
-- Schema: academic
-- Asignaturas por escuela y unidad académica
-- ============================================================

CREATE TABLE academic.subjects (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    school_id uuid NOT NULL,
    academic_unit_id uuid,
    name character varying(255) NOT NULL,
    code character varying(50),
    description text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz,
    CONSTRAINT subjects_pkey PRIMARY KEY (id),
    -- Intra-schema FKs
    CONSTRAINT subjects_school_fkey FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE,
    CONSTRAINT subjects_unit_fkey FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE SET NULL
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.subjects
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
