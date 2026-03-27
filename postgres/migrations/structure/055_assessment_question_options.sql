-- ============================================================
-- 055: assessment.question_options
-- Schema: assessment
-- Opciones de respuesta para preguntas de seleccion multiple
-- ============================================================

CREATE TABLE assessment.question_options (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    question_id uuid NOT NULL,
    option_text text NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT question_options_pkey PRIMARY KEY (id),
    -- Intra-schema FK
    CONSTRAINT question_options_question_fkey FOREIGN KEY (question_id) REFERENCES assessment.questions(id) ON DELETE CASCADE
);

CREATE INDEX idx_options_question ON assessment.question_options(question_id, sort_order);
