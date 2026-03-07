-- ============================================================
-- 037: academic.school_concepts
-- Schema: academic
-- Términos activos por institución
-- ============================================================

CREATE TABLE academic.school_concepts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id   UUID NOT NULL,
    term_key    VARCHAR(100) NOT NULL,
    term_value  VARCHAR(200) NOT NULL,
    category    VARCHAR(50) NOT NULL DEFAULT 'general',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT school_concepts_school_key_unique UNIQUE (school_id, term_key),
    CONSTRAINT fk_school_concepts_school FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE
);

CREATE INDEX idx_school_concepts_school ON academic.school_concepts USING btree (school_id);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.school_concepts
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
