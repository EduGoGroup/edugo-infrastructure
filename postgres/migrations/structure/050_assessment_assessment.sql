-- ============================================================
-- 050: assessment.assessment
-- Schema: assessment
-- Evaluaciones generadas a partir de materiales
-- Cross-schema FK (material_id -> content.materials) va en 070
-- NOTA: Se elimina total_questions y sync_questions_count trigger
-- ============================================================

CREATE TABLE assessment.assessment (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    material_id uuid NOT NULL,
    mongo_document_id character varying(24) NOT NULL,
    questions_count integer DEFAULT 0 NOT NULL,
    status character varying(50) DEFAULT 'generated' NOT NULL,
    title character varying(255),
    pass_threshold integer DEFAULT 70,
    max_attempts integer,
    time_limit_minutes integer,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz,
    CONSTRAINT assessment_pkey PRIMARY KEY (id),
    CONSTRAINT assessment_mongo_unique UNIQUE (mongo_document_id),
    CONSTRAINT assessment_pass_threshold_check CHECK (pass_threshold >= 0 AND pass_threshold <= 100),
    CONSTRAINT assessment_status_check CHECK (status IN ('draft', 'generated', 'published', 'archived', 'closed'))
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON assessment.assessment
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
