-- Mock Data: Membresías (relación usuario-escuela-curso)

-- Teacher en Escuela Primaria
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, enrolled_at, created_at, updated_at) VALUES
('e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'teacher', true, NOW(), NOW(), NOW()),
('e2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'teacher', true, NOW(), NOW(), NOW());

-- Students en Primer Grado A
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, enrolled_at, created_at, updated_at) VALUES
('e3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'student', true, NOW(), NOW(), NOW()),
('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'student', true, NOW(), NOW(), NOW());

-- Student en Primer Grado B
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, enrolled_at, created_at, updated_at) VALUES
('e5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'student', true, NOW(), NOW(), NOW());
