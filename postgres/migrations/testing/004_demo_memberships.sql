-- Mock Data: Membresías (relación usuario-escuela-curso)

-- Teacher en Escuela Primaria
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, status, joined_at, created_at, updated_at) VALUES
('mbr_teacher_math_g1a', 'usr_teacher_math', 'sch_demo_primary', 'au_primary_g1_a', 'teacher', 'active', NOW(), NOW(), NOW()),
('mbr_teacher_science_g1b', 'usr_teacher_science', 'sch_demo_primary', 'au_primary_g1_b', 'teacher', 'active', NOW(), NOW(), NOW());

-- Students en Primer Grado A
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, status, joined_at, created_at, updated_at) VALUES
('mbr_student1_g1a', 'usr_student_1', 'sch_demo_primary', 'au_primary_g1_a', 'student', 'active', NOW(), NOW(), NOW()),
('mbr_student2_g1a', 'usr_student_2', 'sch_demo_primary', 'au_primary_g1_a', 'student', 'active', NOW(), NOW(), NOW());

-- Student en Primer Grado B
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, status, joined_at, created_at, updated_at) VALUES
('mbr_student3_g1b', 'usr_student_3', 'sch_demo_primary', 'au_primary_g1_b', 'student', 'active', NOW(), NOW(), NOW());
