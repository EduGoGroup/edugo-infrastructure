-- ============================================================
-- 052: assessment.assessment_attempt_answer
-- Schema: assessment
-- Respuestas individuales por pregunta en un intento de evaluación
-- ============================================================

CREATE TABLE assessment.assessment_attempt_answer (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    attempt_id uuid NOT NULL,
    question_index integer NOT NULL,
    question_id uuid,
    student_answer text,
    review_status character varying(20) DEFAULT 'pending',
    is_correct boolean,
    points_earned numeric(5,2),
    max_points numeric(5,2),
    time_spent_seconds integer,
    answered_at timestamptz DEFAULT now() NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT assessment_attempt_answer_pkey PRIMARY KEY (id),
    CONSTRAINT assessment_attempt_answer_unique_question UNIQUE (attempt_id, question_index),
    CONSTRAINT assessment_attempt_answer_time_spent_seconds_check CHECK (time_spent_seconds >= 0),
    CONSTRAINT assessment_attempt_answer_review_status_check CHECK (review_status IN ('pending', 'auto_graded', 'reviewed')),
    -- Intra-schema FK
    CONSTRAINT assessment_attempt_answer_attempt_fkey FOREIGN KEY (attempt_id) REFERENCES assessment.assessment_attempt(id) ON DELETE CASCADE
    -- question_id FK va en 070 (assessment.questions se crea en 054, despues de esta tabla)
);
