-- ============================================================
-- 093: academic.announcements
-- Schema: academic
-- Comunicaciones escuela-familia
-- ============================================================

CREATE TABLE academic.announcements (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    school_id uuid NOT NULL,
    academic_unit_id uuid,
    author_id uuid NOT NULL,
    title character varying(200) NOT NULL,
    body text NOT NULL,
    scope character varying(20) NOT NULL,
    target_roles text[] DEFAULT '{}',
    is_pinned boolean DEFAULT false NOT NULL,
    published_at timestamptz DEFAULT now(),
    expires_at timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT announcements_pkey PRIMARY KEY (id),
    CONSTRAINT announcements_scope_check CHECK (scope IN ('school', 'unit', 'role')),
    CONSTRAINT announcements_school_fkey FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE,
    CONSTRAINT announcements_unit_fkey FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE CASCADE
);

CREATE INDEX idx_announcements_school ON academic.announcements(school_id);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.announcements
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
