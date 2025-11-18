-- Mock Data: Unidades académicas de demostración

-- Cursos de la Escuela Primaria
INSERT INTO academic_units (id, school_id, name, unit_type, parent_id, created_at, updated_at) VALUES
-- Grados
('au_primary_grade1', 'sch_demo_primary', 'Primer Grado', 'grade', NULL, NOW(), NOW()),
('au_primary_grade2', 'sch_demo_primary', 'Segundo Grado', 'grade', NULL, NOW(), NOW()),
('au_primary_grade3', 'sch_demo_primary', 'Tercer Grado', 'grade', NULL, NOW(), NOW()),

-- Secciones de Primer Grado
('au_primary_g1_a', 'sch_demo_primary', 'Sección A', 'section', 'au_primary_grade1', NOW(), NOW()),
('au_primary_g1_b', 'sch_demo_primary', 'Sección B', 'section', 'au_primary_grade1', NOW(), NOW());

-- Cursos del Colegio Secundario
INSERT INTO academic_units (id, school_id, name, unit_type, parent_id, created_at, updated_at) VALUES
-- Años
('au_secondary_year1', 'sch_demo_secondary', 'Primer Año', 'year', NULL, NOW(), NOW()),
('au_secondary_year2', 'sch_demo_secondary', 'Segundo Año', 'year', NULL, NOW(), NOW()),

-- Divisiones de Primer Año
('au_secondary_y1_div1', 'sch_demo_secondary', 'División 1', 'division', 'au_secondary_year1', NOW(), NOW()),
('au_secondary_y1_div2', 'sch_demo_secondary', 'División 2', 'division', 'au_secondary_year1', NOW(), NOW());

-- Cursos del Instituto Técnico
INSERT INTO academic_units (id, school_id, name, unit_type, parent_id, created_at, updated_at) VALUES
('au_tech_programming', 'sch_demo_tech', 'Programación I', 'course', NULL, NOW(), NOW()),
('au_tech_databases', 'sch_demo_tech', 'Bases de Datos', 'course', NULL, NOW(), NOW());
