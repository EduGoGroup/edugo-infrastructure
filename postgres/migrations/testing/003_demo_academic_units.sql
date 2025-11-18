-- Mock Data: Unidades académicas de demostración

-- Cursos de la Escuela Primaria
INSERT INTO academic_units (id, school_id, name, code, type, parent_unit_id, created_at, updated_at) VALUES
-- Grados
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Primer Grado', 'G1', 'grade', NULL, NOW(), NOW()),
('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Segundo Grado', 'G2', 'grade', NULL, NOW(), NOW()),
('c3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Tercer Grado', 'G3', 'grade', NULL, NOW(), NOW()),

-- Secciones de Primer Grado
('c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Sección A', 'G1A', 'section', 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), NOW()),
('c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Sección B', 'G1B', 'section', 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), NOW());

-- Cursos del Colegio Secundario
INSERT INTO academic_units (id, school_id, name, code, type, parent_unit_id, created_at, updated_at) VALUES
-- Grados (Secundaria)
('c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Primer Año', 'Y1', 'grade', NULL, NOW(), NOW()),
('c7eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Segundo Año', 'Y2', 'grade', NULL, NOW(), NOW()),

-- Secciones de Primer Año
('c8eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Sección 1', 'Y1S1', 'section', 'c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', NOW(), NOW()),
('c9eebc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Sección 2', 'Y1S2', 'section', 'c6eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', NOW(), NOW());

-- Cursos del Instituto Técnico (Departamentos)
INSERT INTO academic_units (id, school_id, name, code, type, parent_unit_id, created_at, updated_at) VALUES
('d1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Programación I', 'PROG1', 'department', NULL, NOW(), NOW()),
('d2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Bases de Datos', 'DB1', 'department', NULL, NOW(), NOW());
