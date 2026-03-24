-- ============================================================
-- 038: academic.membership_subjects
-- Schema: academic
-- Tabla de union entre memberships y subjects (reemplaza metadata.subjects)
-- ============================================================

CREATE TABLE academic.membership_subjects (
    membership_id uuid NOT NULL,
    subject_id uuid NOT NULL,
    CONSTRAINT membership_subjects_pkey PRIMARY KEY (membership_id, subject_id),
    CONSTRAINT membership_subjects_membership_fkey FOREIGN KEY (membership_id) REFERENCES academic.memberships(id) ON DELETE CASCADE,
    CONSTRAINT membership_subjects_subject_fkey FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE CASCADE
);

CREATE INDEX idx_membership_subjects_membership ON academic.membership_subjects(membership_id);
CREATE INDEX idx_membership_subjects_subject ON academic.membership_subjects(subject_id);
