-- ============================================================
-- 031: academic.academic_units
-- Schema: academic
-- Unidades académicas jerárquicas (grados, cursos, secciones, etc.)
-- ============================================================

CREATE TABLE academic.academic_units (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    parent_unit_id uuid,
    school_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(50) NOT NULL,
    type character varying(50) NOT NULL,
    description text,
    level character varying(50),
    academic_year integer DEFAULT 0,
    metadata jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz,
    CONSTRAINT academic_units_pkey PRIMARY KEY (id),
    CONSTRAINT academic_units_unique_code UNIQUE (school_id, code, academic_year),
    CONSTRAINT academic_units_no_self_reference CHECK (id <> parent_unit_id),
    CONSTRAINT academic_units_type_check CHECK (type IN ('school', 'grade', 'class', 'section', 'club', 'department')),
    -- Intra-schema FKs
    CONSTRAINT academic_units_parent_fkey FOREIGN KEY (parent_unit_id) REFERENCES academic.academic_units(id) ON DELETE SET NULL,
    CONSTRAINT academic_units_school_fkey FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.academic_units
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

CREATE TRIGGER trigger_prevent_academic_unit_cycles BEFORE INSERT OR UPDATE OF parent_unit_id ON academic.academic_units
    FOR EACH ROW EXECUTE FUNCTION public.prevent_academic_unit_cycles();
