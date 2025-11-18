-- Constraints para tabla assessment

ALTER TABLE assessment ADD CONSTRAINT assessment_material_fkey 
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE;

ALTER TABLE assessment ADD CONSTRAINT assessment_mongo_unique 
    UNIQUE (mongo_document_id);

ALTER TABLE assessment ADD CONSTRAINT assessment_status_check 
    CHECK (status IN ('draft', 'generated', 'published', 'archived', 'closed'));

ALTER TABLE assessment ADD CONSTRAINT assessment_pass_threshold_check 
    CHECK (pass_threshold >= 0 AND pass_threshold <= 100);

-- Trigger para mantener sincronizado questions_count y total_questions
CREATE OR REPLACE FUNCTION sync_questions_count()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.total_questions IS NOT NULL THEN
        NEW.questions_count := NEW.total_questions;
    ELSIF NEW.questions_count IS NOT NULL THEN
        NEW.total_questions := NEW.questions_count;
    ELSE
        NEW.total_questions := 0;
        NEW.questions_count := 0;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_sync_questions_count
    BEFORE INSERT OR UPDATE ON assessment
    FOR EACH ROW
    EXECUTE FUNCTION sync_questions_count();

COMMENT ON TRIGGER trg_sync_questions_count ON assessment IS 'Mantiene sincronizado questions_count y total_questions durante transición';
