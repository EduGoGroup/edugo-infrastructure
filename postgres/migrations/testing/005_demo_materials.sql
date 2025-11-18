-- Mock Data: Materiales educativos de demostración

INSERT INTO materials (id, school_id, uploaded_by_teacher_id, academic_unit_id, title, description, file_url, file_type, file_size_bytes, created_at, updated_at) VALUES
-- Material de Matemáticas
('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'Guía de Sumas', 'Material educativo sobre sumas básicas', 'https://s3.example.com/materials/math/suma.pdf', 'application/pdf', 1048576, NOW(), NOW()),
('f2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'Guía de Restas', 'Material educativo sobre restas básicas', 'https://s3.example.com/materials/math/resta.pdf', 'application/pdf', 950000, NOW(), NOW()),

-- Material de Ciencias
('f3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'Las Plantas', 'Video educativo sobre plantas', 'https://s3.example.com/materials/science/plantas.mp4', 'video/mp4', 52428800, NOW(), NOW()),
('f4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'El Ciclo del Agua', 'Presentación sobre el ciclo del agua', 'https://s3.example.com/materials/science/agua.pptx', 'application/vnd.ms-powerpoint', 2097152, NOW(), NOW());
