-- ====================================================================
-- SEEDS: Tipos de concepto y definiciones de terminologia
-- Idempotente: usa ON CONFLICT DO NOTHING
-- ====================================================================

BEGIN;

-- 5 tipos de institucion
INSERT INTO academic.concept_types (id, name, code, description) VALUES
    ('c1000000-0000-0000-0000-000000000001', 'Escuela Primaria',    'primary_school',    'Institucion de educacion basica'),
    ('c1000000-0000-0000-0000-000000000002', 'Colegio Secundario',  'high_school',       'Institucion de educacion media'),
    ('c1000000-0000-0000-0000-000000000003', 'Academia de Idiomas', 'language_academy',  'Centro de ensenanza de idiomas'),
    ('c1000000-0000-0000-0000-000000000004', 'Instituto Tecnico',   'technical_school',  'Formacion tecnica y profesional'),
    ('c1000000-0000-0000-0000-000000000005', 'Taller / Workshop',   'workshop',          'Cursos cortos y talleres practicos')
ON CONFLICT (id) DO NOTHING;

-- Definiciones para "Escuela Primaria"
INSERT INTO academic.concept_definitions (concept_type_id, term_key, term_value, category, sort_order) VALUES
    ('c1000000-0000-0000-0000-000000000001', 'org.name_singular',     'Escuela',       'org',     1),
    ('c1000000-0000-0000-0000-000000000001', 'org.name_plural',       'Escuelas',      'org',     2),
    ('c1000000-0000-0000-0000-000000000001', 'unit.level1',           'Grado',         'unit',    3),
    ('c1000000-0000-0000-0000-000000000001', 'unit.level2',           'Clase',         'unit',    4),
    ('c1000000-0000-0000-0000-000000000001', 'unit.period',           'Periodo',       'unit',    5),
    ('c1000000-0000-0000-0000-000000000001', 'member.student',        'Estudiante',    'member',  6),
    ('c1000000-0000-0000-0000-000000000001', 'member.teacher',        'Profesor',      'member',  7),
    ('c1000000-0000-0000-0000-000000000001', 'member.guardian',       'Acudiente',     'member',  8),
    ('c1000000-0000-0000-0000-000000000001', 'content.subject',       'Materia',       'content', 9),
    ('c1000000-0000-0000-0000-000000000001', 'content.assessment',    'Evaluacion',    'content', 10)
ON CONFLICT (concept_type_id, term_key) DO NOTHING;

-- Definiciones para "Colegio Secundario"
INSERT INTO academic.concept_definitions (concept_type_id, term_key, term_value, category, sort_order) VALUES
    ('c1000000-0000-0000-0000-000000000002', 'org.name_singular',     'Colegio',       'org',     1),
    ('c1000000-0000-0000-0000-000000000002', 'org.name_plural',       'Colegios',      'org',     2),
    ('c1000000-0000-0000-0000-000000000002', 'unit.level1',           'Ano',           'unit',    3),
    ('c1000000-0000-0000-0000-000000000002', 'unit.level2',           'Division',      'unit',    4),
    ('c1000000-0000-0000-0000-000000000002', 'unit.period',           'Trimestre',     'unit',    5),
    ('c1000000-0000-0000-0000-000000000002', 'member.student',        'Alumno',        'member',  6),
    ('c1000000-0000-0000-0000-000000000002', 'member.teacher',        'Docente',       'member',  7),
    ('c1000000-0000-0000-0000-000000000002', 'member.guardian',       'Tutor',         'member',  8),
    ('c1000000-0000-0000-0000-000000000002', 'content.subject',       'Asignatura',    'content', 9),
    ('c1000000-0000-0000-0000-000000000002', 'content.assessment',    'Examen',        'content', 10)
ON CONFLICT (concept_type_id, term_key) DO NOTHING;

-- Definiciones para "Academia de Idiomas"
INSERT INTO academic.concept_definitions (concept_type_id, term_key, term_value, category, sort_order) VALUES
    ('c1000000-0000-0000-0000-000000000003', 'org.name_singular',     'Academia',      'org',     1),
    ('c1000000-0000-0000-0000-000000000003', 'org.name_plural',       'Academias',     'org',     2),
    ('c1000000-0000-0000-0000-000000000003', 'unit.level1',           'Level',         'unit',    3),
    ('c1000000-0000-0000-0000-000000000003', 'unit.level2',           'Class',         'unit',    4),
    ('c1000000-0000-0000-0000-000000000003', 'unit.period',           'Term',          'unit',    5),
    ('c1000000-0000-0000-0000-000000000003', 'member.student',        'Student',       'member',  6),
    ('c1000000-0000-0000-0000-000000000003', 'member.teacher',        'Teacher',       'member',  7),
    ('c1000000-0000-0000-0000-000000000003', 'member.guardian',       'Parent',        'member',  8),
    ('c1000000-0000-0000-0000-000000000003', 'content.subject',       'Course',        'content', 9),
    ('c1000000-0000-0000-0000-000000000003', 'content.assessment',    'Test',          'content', 10)
ON CONFLICT (concept_type_id, term_key) DO NOTHING;

-- Definiciones para "Instituto Tecnico"
INSERT INTO academic.concept_definitions (concept_type_id, term_key, term_value, category, sort_order) VALUES
    ('c1000000-0000-0000-0000-000000000004', 'org.name_singular',     'Instituto',     'org',     1),
    ('c1000000-0000-0000-0000-000000000004', 'org.name_plural',       'Institutos',    'org',     2),
    ('c1000000-0000-0000-0000-000000000004', 'unit.level1',           'Semestre',      'unit',    3),
    ('c1000000-0000-0000-0000-000000000004', 'unit.level2',           'Seccion',       'unit',    4),
    ('c1000000-0000-0000-0000-000000000004', 'unit.period',           'Cuatrimestre',  'unit',    5),
    ('c1000000-0000-0000-0000-000000000004', 'member.student',        'Aprendiz',      'member',  6),
    ('c1000000-0000-0000-0000-000000000004', 'member.teacher',        'Instructor',    'member',  7),
    ('c1000000-0000-0000-0000-000000000004', 'member.guardian',       'Representante', 'member',  8),
    ('c1000000-0000-0000-0000-000000000004', 'content.subject',       'Modulo',        'content', 9),
    ('c1000000-0000-0000-0000-000000000004', 'content.assessment',    'Prueba',        'content', 10)
ON CONFLICT (concept_type_id, term_key) DO NOTHING;

-- Definiciones para "Taller / Workshop"
INSERT INTO academic.concept_definitions (concept_type_id, term_key, term_value, category, sort_order) VALUES
    ('c1000000-0000-0000-0000-000000000005', 'org.name_singular',     'Taller',        'org',     1),
    ('c1000000-0000-0000-0000-000000000005', 'org.name_plural',       'Talleres',      'org',     2),
    ('c1000000-0000-0000-0000-000000000005', 'unit.level1',           'Modulo',        'unit',    3),
    ('c1000000-0000-0000-0000-000000000005', 'unit.level2',           'Grupo',         'unit',    4),
    ('c1000000-0000-0000-0000-000000000005', 'unit.period',           'Ciclo',         'unit',    5),
    ('c1000000-0000-0000-0000-000000000005', 'member.student',        'Participante',  'member',  6),
    ('c1000000-0000-0000-0000-000000000005', 'member.teacher',        'Facilitador',   'member',  7),
    ('c1000000-0000-0000-0000-000000000005', 'member.guardian',       'Responsable',   'member',  8),
    ('c1000000-0000-0000-0000-000000000005', 'content.subject',       'Taller',        'content', 9),
    ('c1000000-0000-0000-0000-000000000005', 'content.assessment',    'Ejercicio',     'content', 10)
ON CONFLICT (concept_type_id, term_key) DO NOTHING;

COMMIT;
