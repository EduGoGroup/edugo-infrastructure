-- ============================================================
-- 032: academic.memberships
-- Schema: academic
-- Membresías de usuarios en escuelas y unidades académicas
-- Cross-schema FK (user_id -> auth.users) va en 070
-- ============================================================

CREATE TABLE academic.memberships (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    school_id uuid NOT NULL,
    academic_unit_id uuid,
    role character varying(50) NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true NOT NULL,
    enrolled_at timestamptz DEFAULT now() NOT NULL,
    withdrawn_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT memberships_pkey PRIMARY KEY (id),
    CONSTRAINT memberships_unique_membership UNIQUE (user_id, school_id, academic_unit_id, role),
    CONSTRAINT memberships_role_check CHECK (role IN ('teacher', 'student', 'guardian', 'coordinator', 'admin', 'assistant')),
    -- Intra-schema FKs
    CONSTRAINT memberships_school_fkey FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE,
    CONSTRAINT memberships_unit_fkey FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE CASCADE
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.memberships
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
