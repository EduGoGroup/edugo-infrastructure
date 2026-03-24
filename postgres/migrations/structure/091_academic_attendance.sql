-- ============================================================
-- 091: academic.attendance
-- Schema: academic
-- Registro de asistencia de estudiantes
-- ============================================================

CREATE TABLE academic.attendance (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    membership_id uuid NOT NULL,
    subject_id uuid,
    date date NOT NULL,
    status character varying(20) NOT NULL,
    recorded_by uuid NOT NULL,
    notes text,
    created_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT attendance_pkey PRIMARY KEY (id),
    CONSTRAINT attendance_unique UNIQUE (membership_id, subject_id, date),
    CONSTRAINT attendance_status_check CHECK (status IN ('present', 'absent', 'late', 'excused', 'remote')),
    CONSTRAINT attendance_membership_fkey FOREIGN KEY (membership_id) REFERENCES academic.memberships(id) ON DELETE CASCADE,
    CONSTRAINT attendance_subject_fkey FOREIGN KEY (subject_id) REFERENCES academic.subjects(id) ON DELETE CASCADE
);

CREATE INDEX idx_attendance_membership ON academic.attendance(membership_id);
CREATE INDEX idx_attendance_date ON academic.attendance(date);
