-- ============================================================
-- 059: assessment.assessment_stats
-- Schema: assessment
-- Cumulative statistics per assessment, updated by the worker
-- ============================================================

CREATE TABLE assessment.assessment_stats (
    id uuid NOT NULL,
    assessment_id uuid NOT NULL,
    attempt_count integer DEFAULT 0 NOT NULL,
    average_score numeric(5,2) DEFAULT 0 NOT NULL,
    average_percentage numeric(5,2) DEFAULT 0 NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT assessment_stats_pkey PRIMARY KEY (id),
    CONSTRAINT assessment_stats_assessment_unique UNIQUE (assessment_id),
    -- Intra-schema FK
    CONSTRAINT assessment_stats_assessment_fkey FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE
);

-- Note: assessment_id already has a unique index from the UNIQUE constraint.
-- No separate index needed.
