-- Constraints para tabla memberships

ALTER TABLE memberships ADD CONSTRAINT memberships_user_fkey 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE memberships ADD CONSTRAINT memberships_school_fkey 
    FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;

ALTER TABLE memberships ADD CONSTRAINT memberships_unit_fkey 
    FOREIGN KEY (academic_unit_id) REFERENCES academic_units(id) ON DELETE CASCADE;

ALTER TABLE memberships ADD CONSTRAINT memberships_unique_membership 
    UNIQUE(user_id, school_id, academic_unit_id, role);

ALTER TABLE memberships ADD CONSTRAINT memberships_role_check 
    CHECK (role IN ('teacher', 'student', 'guardian', 'coordinator', 'admin', 'assistant'));
