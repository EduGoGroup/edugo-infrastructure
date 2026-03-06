-- ============================================================
-- 035: academic.concept_types
-- Schema: academic
-- Tipos de institucion con terminologia predefinida
-- ============================================================

CREATE TABLE academic.concept_types (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    code            VARCHAR(50) NOT NULL UNIQUE,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_concept_types_code ON academic.concept_types USING btree (code);
CREATE INDEX idx_concept_types_active ON academic.concept_types USING btree (is_active) WHERE (is_active = true);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.concept_types
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
