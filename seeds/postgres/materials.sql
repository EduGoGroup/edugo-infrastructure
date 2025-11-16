-- Seeds de materiales de prueba
-- Ejecutar después de migración 005_create_materials

INSERT INTO materials (id, school_id, uploaded_by_teacher_id, title, description, subject, grade, file_url, file_type, file_size_bytes, status) VALUES
-- Material 1
('66666666-6666-6666-6666-666666666666', 
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222222',
 'Introducción a la Física Cuántica',
 'Material educativo sobre conceptos básicos de física cuántica',
 'Física',
 '10th',
 's3://edugo-materials-dev/fisica-cuantica.pdf',
 'application/pdf',
 2048000,
 'ready'),

-- Material 2
('77777777-7777-7777-7777-777777777777',
 '44444444-4444-4444-4444-444444444444',
 '22222222-2222-2222-2222-222222222222',
 'Álgebra Lineal - Matrices',
 'Ejercicios y teoría sobre matrices y determinantes',
 'Matemáticas',
 '11th',
 's3://edugo-materials-dev/algebra-matrices.pdf',
 'application/pdf',
 1524000,
 'ready'),

-- Material 3 (en procesamiento)
('88888888-8888-8888-8888-888888888888',
 '55555555-5555-5555-5555-555555555555',
 '22222222-2222-2222-2222-222222222222',
 'Historia de Chile - Siglo XX',
 'Material sobre eventos históricos del siglo XX en Chile',
 'Historia',
 '9th',
 's3://edugo-materials-dev/historia-chile.pdf',
 'application/pdf',
 3072000,
 'processing')

ON CONFLICT (id) DO NOTHING;
