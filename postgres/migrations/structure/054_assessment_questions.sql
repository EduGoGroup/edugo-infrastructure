-- ============================================================
-- 054: assessment.questions
-- Schema: assessment
-- Preguntas de evaluacion con soporte para diferentes tipos
-- ============================================================

CREATE TABLE assessment.questions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    assessment_id uuid NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    question_text text NOT NULL,
    question_type character varying(50) NOT NULL,
    correct_answer text,
    explanation text,
    points numeric(5,2) DEFAULT 1 NOT NULL,
    difficulty character varying(20),
    tags text[],
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT questions_pkey PRIMARY KEY (id),
    CONSTRAINT questions_question_type_check CHECK (question_type IN ('multiple_choice', 'true_false', 'short_answer', 'open_ended')),
    CONSTRAINT questions_difficulty_check CHECK (difficulty IN ('easy', 'medium', 'hard')),
    -- Intra-schema FK
    CONSTRAINT questions_assessment_fkey FOREIGN KEY (assessment_id) REFERENCES assessment.assessment(id) ON DELETE CASCADE
);

CREATE INDEX idx_questions_assessment ON assessment.questions(assessment_id, sort_order);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON assessment.questions
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
