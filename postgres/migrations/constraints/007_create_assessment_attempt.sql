-- Constraints para tabla assessment_attempt

ALTER TABLE assessment_attempt ADD CONSTRAINT assessment_attempt_assessment_fkey 
    FOREIGN KEY (assessment_id) REFERENCES assessment(id) ON DELETE CASCADE;

ALTER TABLE assessment_attempt ADD CONSTRAINT assessment_attempt_student_fkey 
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE assessment_attempt ADD CONSTRAINT assessment_attempt_status_check 
    CHECK (status IN ('in_progress', 'completed', 'abandoned'));

ALTER TABLE assessment_attempt ADD CONSTRAINT assessment_attempt_time_spent_seconds_check 
    CHECK (time_spent_seconds IS NULL OR (time_spent_seconds > 0 AND time_spent_seconds <= 7200));

ALTER TABLE assessment_attempt ADD CONSTRAINT unique_idempotency_key 
    UNIQUE (idempotency_key);

ALTER TABLE assessment_attempt ADD CONSTRAINT check_attempt_time_logical 
    CHECK (completed_at IS NULL OR completed_at > started_at);

COMMENT ON CONSTRAINT check_attempt_time_logical ON assessment_attempt IS 'Validar que completed_at > started_at';
