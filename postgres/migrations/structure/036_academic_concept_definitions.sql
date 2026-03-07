-- ============================================================
-- 036: academic.concept_definitions
-- Schema: academic
-- Terminos predeterminados por tipo de institucion (plantilla)
-- ============================================================

CREATE TABLE academic.concept_definitions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    concept_type_id   UUID NOT NULL,
    term_key          VARCHAR(100) NOT NULL,
    term_value        VARCHAR(200) NOT NULL,
    category          VARCHAR(50) NOT NULL DEFAULT 'general',
    sort_order        INT NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT concept_definitions_type_key_unique UNIQUE (concept_type_id, term_key),
    CONSTRAINT fk_concept_definitions_type FOREIGN KEY (concept_type_id) REFERENCES academic.concept_types(id) ON DELETE CASCADE
);

CREATE INDEX idx_concept_definitions_type ON academic.concept_definitions USING btree (concept_type_id);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.concept_definitions
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();
