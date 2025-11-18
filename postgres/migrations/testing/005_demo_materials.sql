-- Mock Data: Materiales educativos de demostración

INSERT INTO materials (id, school_id, academic_unit_id, title, description, content_type, storage_provider, storage_key, file_size, mime_type, created_by, created_at, updated_at) VALUES
-- Material de Matemáticas
('mat_math_suma', 'sch_demo_primary', 'au_primary_g1_a', 'Guía de Sumas', 'Material educativo sobre sumas básicas', 'document', 's3', 'materials/math/suma.pdf', 1048576, 'application/pdf', 'usr_teacher_math', NOW(), NOW()),
('mat_math_resta', 'sch_demo_primary', 'au_primary_g1_a', 'Guía de Restas', 'Material educativo sobre restas básicas', 'document', 's3', 'materials/math/resta.pdf', 950000, 'application/pdf', 'usr_teacher_math', NOW(), NOW()),

-- Material de Ciencias
('mat_science_plantas', 'sch_demo_primary', 'au_primary_g1_b', 'Las Plantas', 'Video educativo sobre plantas', 'video', 's3', 'materials/science/plantas.mp4', 52428800, 'video/mp4', 'usr_teacher_science', NOW(), NOW()),
('mat_science_agua', 'sch_demo_primary', 'au_primary_g1_b', 'El Ciclo del Agua', 'Presentación sobre el ciclo del agua', 'presentation', 's3', 'materials/science/agua.pptx', 2097152, 'application/vnd.ms-powerpoint', 'usr_teacher_science', NOW(), NOW());
