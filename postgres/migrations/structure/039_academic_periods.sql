-- ============================================================
-- 039: academic.academic_periods
-- Schema: academic
-- Periodos academicos (semestre, trimestre, bimestre)
-- ============================================================

CREATE TABLE academic.academic_periods (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    school_id uuid NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(20),
    type character varying(20) NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    is_active boolean DEFAULT false NOT NULL,
    academic_year integer NOT NULL,
    sort_order integer DEFAULT 0,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT academic_periods_pkey PRIMARY KEY (id),
    CONSTRAINT academic_periods_type_check CHECK (type IN ('semester', 'trimester', 'bimester', 'quarter')),
    CONSTRAINT academic_periods_school_fkey FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE
);

CREATE INDEX idx_academic_periods_school ON academic.academic_periods(school_id);
CREATE UNIQUE INDEX idx_academic_periods_active ON academic.academic_periods(school_id) WHERE is_active = true;

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.academic_periods
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
