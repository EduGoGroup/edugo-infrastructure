-- ============================================================
-- 050: assessment.assessment
-- Schema: assessment
-- Evaluaciones con soporte N:N a materiales (via 053)
-- Cross-schema FKs (school_id, created_by_user_id) van en 070
-- ============================================================

CREATE TABLE assessment.assessment (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    mongo_document_id character varying(24),
    source_type character varying(20) DEFAULT 'manual' NOT NULL,
    school_id uuid,
    created_by_user_id uuid,
    questions_count integer DEFAULT 0 NOT NULL,
    status character varying(50) DEFAULT 'generated' NOT NULL,
    title character varying(255),
    description text,
    pass_threshold numeric(5,2) DEFAULT 70,
    max_attempts integer,
    time_limit_minutes numeric(7,2),
    is_timed boolean NOT NULL DEFAULT false,
    shuffle_questions boolean NOT NULL DEFAULT false,
    show_correct_answers boolean NOT NULL DEFAULT true,
    available_from timestamptz,
    available_until timestamptz,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz,
    CONSTRAINT assessment_pkey PRIMARY KEY (id),
    CONSTRAINT assessment_mongo_unique UNIQUE (mongo_document_id),
    CONSTRAINT assessment_pass_threshold_check CHECK (pass_threshold >= 0 AND pass_threshold <= 100),
    CONSTRAINT assessment_status_check CHECK (status IN ('draft', 'generated', 'published', 'archived', 'closed')),
    CONSTRAINT assessment_available_dates_check CHECK (available_until IS NULL OR available_from IS NULL OR available_until > available_from),
    CONSTRAINT assessment_source_type_check CHECK (source_type IN ('manual', 'ai_generated'))
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON assessment.assessment
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
