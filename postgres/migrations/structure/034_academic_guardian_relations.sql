-- ============================================================
-- 034: academic.guardian_relations
-- Schema: academic
-- Relaciones apoderado-estudiante
-- Cross-schema FKs (guardian_id, student_id -> auth.users) van en 070
-- ============================================================

CREATE TABLE academic.guardian_relations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    guardian_id uuid NOT NULL,
    student_id uuid NOT NULL,
    relationship_type character varying(50) DEFAULT 'parent' NOT NULL,
    is_primary boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    status character varying(20) NOT NULL DEFAULT 'active',
    created_by uuid,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT guardian_relations_pkey PRIMARY KEY (id),
    CONSTRAINT guardian_relations_unique UNIQUE (guardian_id, student_id),
    CONSTRAINT guardian_relations_type_check CHECK (relationship_type IN ('parent', 'guardian', 'tutor', 'other')),
    CONSTRAINT guardian_relations_status_check CHECK (status IN ('pending', 'active', 'rejected', 'revoked'))
);

CREATE INDEX idx_guardian_relations_status ON academic.guardian_relations (status);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.guardian_relations
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
