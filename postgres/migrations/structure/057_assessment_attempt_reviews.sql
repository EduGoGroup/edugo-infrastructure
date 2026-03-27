-- ============================================================
-- 057: assessment.attempt_reviews
-- Schema: assessment
-- Revisiones de respuestas de intentos por profesores
-- Cross-schema FK (reviewer_id -> auth.users) va en 070
-- ============================================================

CREATE TABLE assessment.attempt_reviews (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    attempt_answer_id uuid NOT NULL,
    reviewer_id uuid NOT NULL,
    points_awarded numeric(5,2) NOT NULL,
    feedback text,
    reviewed_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT attempt_reviews_pkey PRIMARY KEY (id),
    CONSTRAINT unique_review UNIQUE (attempt_answer_id),
    -- Intra-schema FK
    CONSTRAINT attempt_reviews_answer_fkey FOREIGN KEY (attempt_answer_id) REFERENCES assessment.assessment_attempt_answer(id) ON DELETE CASCADE
);
