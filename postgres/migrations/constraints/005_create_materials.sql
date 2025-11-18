-- Constraints para tabla materials

ALTER TABLE materials ADD CONSTRAINT materials_school_fkey 
    FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;

ALTER TABLE materials ADD CONSTRAINT materials_teacher_fkey 
    FOREIGN KEY (uploaded_by_teacher_id) REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE materials ADD CONSTRAINT materials_unit_fkey 
    FOREIGN KEY (academic_unit_id) REFERENCES academic_units(id) ON DELETE SET NULL;

ALTER TABLE materials ADD CONSTRAINT materials_status_check 
    CHECK (status IN ('uploaded', 'processing', 'ready', 'failed'));
