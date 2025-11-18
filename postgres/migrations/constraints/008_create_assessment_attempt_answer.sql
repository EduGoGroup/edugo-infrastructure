-- Constraints para tabla assessment_attempt_answer

ALTER TABLE assessment_attempt_answer ADD CONSTRAINT assessment_attempt_answer_attempt_fkey 
    FOREIGN KEY (attempt_id) REFERENCES assessment_attempt(id) ON DELETE CASCADE;

ALTER TABLE assessment_attempt_answer ADD CONSTRAINT assessment_attempt_answer_unique_question 
    UNIQUE(attempt_id, question_index);

ALTER TABLE assessment_attempt_answer ADD CONSTRAINT assessment_attempt_answer_time_spent_seconds_check 
    CHECK (time_spent_seconds >= 0);
