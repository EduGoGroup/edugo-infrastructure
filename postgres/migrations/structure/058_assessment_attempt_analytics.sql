-- ============================================================
-- 058: assessment.attempt_analytics
-- Schema: assessment
-- Analytics data derived from assessment attempts by the worker
-- ============================================================

CREATE TABLE assessment.attempt_analytics (
    id uuid NOT NULL,
    attempt_id uuid NOT NULL,
    assessment_id uuid NOT NULL,
    student_id uuid NOT NULL,
    school_id uuid NOT NULL,
    score numeric(5,2) NOT NULL,
    total_points numeric(5,2) NOT NULL,
    percentage numeric(5,2),
    duration_seconds integer,
    submitted_at timestamptz NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT attempt_analytics_pkey PRIMARY KEY (id),
    CONSTRAINT attempt_analytics_attempt_unique UNIQUE (attempt_id),
    -- Intra-schema FK
    CONSTRAINT attempt_analytics_attempt_fkey FOREIGN KEY (attempt_id) REFERENCES assessment.assessment_attempt(id) ON DELETE CASCADE,
    CONSTRAINT attempt_analytics_assessment_fkey FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE
);

CREATE INDEX idx_analytics_assessment ON assessment.attempt_analytics(assessment_id);
CREATE INDEX idx_analytics_student ON assessment.attempt_analytics(student_id);
CREATE INDEX idx_analytics_school ON assessment.attempt_analytics(school_id);
CREATE INDEX idx_analytics_submitted ON assessment.attempt_analytics(submitted_at DESC);
