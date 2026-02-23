-- ============================================================
-- 051: assessment.assessment_attempt
-- Schema: assessment
-- Intentos de evaluación por estudiantes
-- Cross-schema FK (student_id -> auth.users) va en 070
-- ============================================================

CREATE TABLE assessment.assessment_attempt (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    assessment_id uuid NOT NULL,
    student_id uuid NOT NULL,
    started_at timestamptz DEFAULT now() NOT NULL,
    completed_at timestamptz,
    score numeric(5,2),
    max_score numeric(5,2),
    percentage numeric(5,2),
    status character varying(50) DEFAULT 'in_progress' NOT NULL,
    time_spent_seconds integer,
    idempotency_key character varying(64),
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT assessment_attempt_pkey PRIMARY KEY (id),
    CONSTRAINT unique_idempotency_key UNIQUE (idempotency_key),
    CONSTRAINT assessment_attempt_status_check CHECK (status IN ('in_progress', 'completed', 'abandoned')),
    CONSTRAINT assessment_attempt_time_spent_seconds_check CHECK (time_spent_seconds IS NULL OR (time_spent_seconds > 0 AND time_spent_seconds <= 7200)),
    CONSTRAINT check_attempt_time_logical CHECK (completed_at IS NULL OR completed_at > started_at),
    -- Intra-schema FK
    CONSTRAINT assessment_attempt_assessment_fkey FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE
);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON assessment.assessment_attempt
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
