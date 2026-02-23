-- ============================================================
-- 040: content.materials
-- Schema: content
-- Materiales educativos subidos por docentes
-- Cross-schema FKs (school_id, uploaded_by_teacher_id, academic_unit_id) van en 070
-- ============================================================

CREATE TABLE content.materials (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    school_id uuid NOT NULL,
    uploaded_by_teacher_id uuid NOT NULL,
    academic_unit_id uuid,
    title character varying(255) NOT NULL,
    description text,
    subject character varying(100),
    grade character varying(50),
    file_url text NOT NULL,
    file_type character varying(100) NOT NULL,
    file_size_bytes bigint NOT NULL,
    status character varying(50) DEFAULT 'uploaded' NOT NULL,
    processing_started_at timestamptz,
    processing_completed_at timestamptz,
    is_public boolean DEFAULT false NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz,
    CONSTRAINT materials_pkey PRIMARY KEY (id),
    CONSTRAINT materials_status_check CHECK (status IN ('uploaded', 'processing', 'ready', 'failed'))
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON content.materials
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
